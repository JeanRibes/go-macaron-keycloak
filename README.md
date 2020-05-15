# Explications
Ce serveur (écrit en Go) utilise une base de données MongoDB et se connecte à Keycloak.

Il récupère les profils adhérents depuis Keycloak quand ceux-ci se connectent en OpenIDconnect.
Si un utilisateur possède le rôle ``"bureau"``, il est administrateur.

Les profils utilisateurs sont mis à jour à chaque connexion. Il n'existe pas d'autre mode de connexion que Keycloak (ou autres Openidconnect).
## Permissions
Mode alternatif où les adhérents n'ont pas besoin de se connecter
On truste les admins d'asso et on leur permet de rechercher par email dans la base de données.
C'est forcé dans l'appli ...

## Fonctionnement
Utilise le *framework* **go-macaron**, qui utilise le même principe de Handlers et middlewares que `net/http`.
Le routage est interopérable avec des Handlers Go `net/http` classiques.

Les templates sont faits avec **pongo2**, qui a la même syntaxe que les templates *Django*

La base de données est **MongoDB**, du *NoSQL*.
# Déploiement
## Configuration
### Variables d'environnment
```shell script
KEYCLOAK_URL="https://sso.asso-insa-lyon.fr/auth/realms/asso-insa-lyon"
CLIENT_ID=""
CLIENT_SECRET="1337-secret-4242-hax0r"
BASE_URL=""
```
### Keycloak
Assigner à un groupe Bureau AMI le rôle client goami/bureau
