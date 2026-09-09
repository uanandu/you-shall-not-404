---
title: Un seul fichier de config, tous les services locaux
description: Notes de conception sur DevRouter — un proxy de développement local et un tunnel auto-hébergé en Go que je construis, et pourquoi un seul fichier YAML devrait battre un dossier de scripts.
date: 2026-01-14
tags: [go, devrouter, outillage]
lang: fr
draft: true
---

Chaque équipe que j'ai rejointe a une version du même problème : une poignée de services qui ont chacun besoin de leur propre port, de leur propre `.env`, et d'un README que personne ne met à jour. Les nouvelles recrues passent leur première journée juste à faire fonctionner `localhost`.

**DevRouter** est ce que je construis pour régler ça. Le plan : un seul fichier de configuration démarre les services, les route par nom d'hôte, et `devrouter share` expose n'importe lequel d'entre eux publiquement via un tunnel WebSocket que vous contrôlez — pas de service de tunnel tiers, pas de facture surprise.

## Pourquoi router par nom d'hôte plutôt que par port

Les ports sont un détail d'implémentation qui s'infiltre dans le flux de travail quotidien. `api.local` et `web.local` sont des adresses ; `:3001` et `:5173` sont des détails triviaux qu'il faut se rappeler. Une fois le routage fait par nom d'hôte, ajouter un nouveau service devient un simple diff de config, pas un changement de contexte.

```yaml
services:
  api:
    host: api.local
    cmd: go run ./cmd/api
  web:
    host: web.local
    cmd: pnpm dev
```

## Le tunnel est la partie difficile

Faire du proxy localement, c'est facile. Exposer un service local en toute sécurité, sans confier son trafic au réseau de bordure de quelqu'un d'autre, ce ne l'est pas. Le plan pour le tunnel WebSocket de DevRouter est de le garder délibérément simple — pas de magie DNS dynamique, pas de système de plugins — parce que la version simple est celle en qui j'aurai assez confiance pour l'utiliser lors de démos client, une fois construite.

Go est le bon choix ici, spécifiquement pour la vitesse d'itération : un client et un serveur de tunnel assez petits pour me permettre de réécrire le protocole de framing plus d'une fois en cherchant à régler la contre-pression, sans qu'aucune des deux réécritures ne prenne une journée complète.
