# 🔄 COMEBACK DATA - Script de Restauration

Ce script permet de récupérer tous les fichiers d'un utilisateur depuis l'API et de les restaurer à leur emplacement original.

## 🚀 Utilisation

### Prérequis
1. Go installé sur votre machine
2. L'API doit être démarrée sur `http://localhost:8080`
3. Une victime doit exister dans l'API avec des fichiers uploadés

### Configuration
Modifiez les constantes dans le script selon vos besoins :
```go
const (
    API_BASE_URL = "http://localhost:8080"  // URL de votre API
    VICTIM_NAME  = "John"                   // Nom de la victime
    PASSWORD     = "TEST123"                // Mot de passe de la victime
)
```

### Exécution
```bash
# Naviguer vers le dossier comeback
cd comeback

# Exécuter le script
go run comeback_data.go
```

## ✨ Fonctionnalités

- ✅ **Récupération automatique** de la liste des fichiers
- ✅ **Restauration à l'emplacement original** (filepathvictime)
- ✅ **Création automatique des dossiers** manquants
- ✅ **Gestion des chemins Windows** (conversion automatique)
- ✅ **Feedback visuel** avec emojis et statistiques
- ✅ **Gestion d'erreurs** robuste
- ✅ **Authentification** sécurisée

## 🔄 Workflow du script

1. **📋** Récupère la liste des fichiers depuis l'API
2. **📥** Télécharge chaque fichier
3. **📁** Crée les dossiers parents si nécessaire
4. **💾** Restaure le fichier à son emplacement original
5. **📊** Affiche les statistiques de restauration

## 📊 Exemple de sortie

```
🔄 COMEBACK DATA - Restauration des fichiers
=============================================
🔄 Début de la restauration des données...
👤 Victime: John
🌐 API: http://localhost:8080
📋 Récupération de la liste des fichiers...
✅ 3 fichiers trouvés pour John

📁 Restauration de 3 fichiers...
📥 Téléchargement: document1.txt
✅ Fichier restauré: C:\Users\John\Documents\document1.txt
📥 Téléchargement: document2.pdf
✅ Fichier restauré: C:\Users\John\Documents\document2.pdf
📥 Téléchargement: image.jpg
✅ Fichier restauré: C:\Users\John\Pictures\image.jpg

🎉 Restauration terminée!
✅ Fichiers restaurés avec succès: 3

✨ Restauration terminée avec succès!
```

## 🔧 Personnalisation

### Changer les identifiants de victime
Modifiez les constantes `VICTIM_NAME` et `PASSWORD` dans le script.

### Changer l'URL de l'API
Modifiez la constante `API_BASE_URL` si votre API n'est pas sur localhost:8080.

### Ajouter des filtres
Vous pouvez modifier le script pour restaurer seulement certains types de fichiers :
```go
// Restaurer seulement les fichiers PDF
if !strings.HasSuffix(file.FileName, ".pdf") {
    continue
}
```

## 🐛 Dépannage

### Erreur "Authentification échouée"
Vérifiez que la victime existe dans l'API avec les bons identifiants.

### Erreur "Aucun fichier à restaurer"
La victime n'a pas de fichiers uploadés dans l'API.

### Erreur de permissions
Vérifiez que vous avez les droits d'écriture dans les dossiers de destination.

### Erreur de connexion
Vérifiez que l'API est démarrée sur le bon port.

## 🔒 Sécurité

- ✅ Authentification requise pour accéder aux fichiers
- ✅ Chaque victime ne peut restaurer que ses propres fichiers
- ✅ Vérification des permissions avant création des dossiers
- ✅ Gestion sécurisée des chemins de fichiers

## 📝 Notes importantes

- Le script utilise le champ `FilePathVictime` pour restaurer à l'emplacement original
- Les dossiers parents sont créés automatiquement si ils n'existent pas
- Les chemins Windows sont convertis automatiquement
- Le script continue même si certains fichiers échouent 