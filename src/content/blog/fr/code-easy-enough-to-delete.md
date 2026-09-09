---
title: Du code assez simple pour être supprimé
description: Sur le fait d'écrire des systèmes pour des clients du secteur public où le vrai livrable est la capacité de partir proprement.
date: 2025-11-08
tags: [devops, consultation, kubernetes]
lang: fr
draft: true
---

Faire de la consultation sur de l'infrastructure gouvernementale enseigne une leçon facile à sauter quand on construit pour soi-même : le meilleur système n'est pas le plus astucieux à construire, c'est celui dont la prochaine équipe peut supprimer une partie sans danger.

En travaillant sur des mandats avec EDSC et le gouvernement du Québec, les systèmes qui ont tenu le coup face au roulement de personnel, aux cycles budgétaires et aux audits n'étaient pas ceux avec les abstractions les plus sophistiquées. C'étaient ceux où un seul chart Helm, une seule application ArgoCD, ou un seul module Terraform pouvait être retiré sans que personne n'ait à faire d'archéologie pour découvrir ce dont il était secrètement porteur.

## Ce que « facile à supprimer » veut vraiment dire

- **Aucune dépendance implicite entre les services.** Si retirer un composant demande une enquête dans un canal Slack pour découvrir ce qui d'autre va casser, ce n'était pas découplé — ça en avait seulement l'air.
- **La configuration vit avec ce qu'elle configure.** Les dépôts de configuration centralisés sont pratiques jusqu'à ce que la personne qui comprend le mapping quitte.
- **Chaque déploiement est reproductible à partir de la source**, pas d'un runbook dont quelqu'un se souvient à moitié. ArgoCD force ça par construction, ce qui explique en grande partie pourquoi j'y ai recours par défaut.

## La partie inconfortable

Ça veut dire résister à l'envie de construire la version élégante et générale de quelque chose quand un client n'a besoin que de la version spécifique aujourd'hui. La généralité est un pari sur le fait que le contexte d'une future équipe correspondra au vôtre — et dans le secteur public, cette équipe est rarement celle assise en face de vous.

Les systèmes dont je suis le plus fier ne sont pas ceux qui en font le plus. Ce sont ceux où, un an plus tard, supprimer le mauvais fichier ne fait pas tomber trois choses sans rapport entre elles.
