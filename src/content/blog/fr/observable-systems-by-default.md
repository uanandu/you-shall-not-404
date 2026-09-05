---
title: Des systèmes observables par défaut
description: Pourquoi je veux arrêter d'ajouter la journalisation après coup, et commencer à concevoir les points de terminaison autour des traces qu'ils devraient produire.
date: 2026-03-02
tags: [devops, observabilité, api]
lang: fr
---

La majorité du travail d'observabilité se fait après qu'un problème soit déjà survenu. Un service se comporte mal, personne ne peut expliquer pourquoi, et le correctif est un `console.log` posé à la hâte ou une intégration Datadog câblée sous pression. Ça fonctionne, mais c'est réactif — le système n'a jamais été conçu pour être observé, il a été conçu puis *interrogé*.

Concevoir **API Darwin** est ce qui change ma façon de penser à ça. Le plan est d'ingérer des spécifications OpenAPI et des journaux de trafic pour construire un graphe de relations entre points de terminaison — ce qui force à poser une question différente dès le départ : pas « comment vais-je déboguer ça plus tard », mais « qu'est-ce que ce point de terminaison doit pouvoir dire de lui-même, dès maintenant ».

## Trois habitudes que je bâtis

1. **Des journaux structurés dès le premier commit.** Pas parce qu'on en aura besoin aujourd'hui, mais parce que restructurer des journaux en texte libre plus tard est franchement pénible.
2. **Nommer les choses pour le graphe, pas pour le fichier.** Le nom d'un point de terminaison devrait décrire sa relation avec le reste du système, pas seulement sa route — le genre de refactoring qui tend à payer le plus sur des bases de code héritées.
3. **Traiter la dépréciation comme un état de première classe**, pas comme un commentaire après coup. Si un système ne peut pas dire ce qui est en train de mourir, il ne peut pas dire ce sur quoi il est sécuritaire de construire.

Rien de tout ça n'est exotique. C'est surtout de la discipline — le genre qu'il est facile de reporter sous une échéance et coûteux d'ignorer une fois qu'un système a de vrais utilisateurs.

## Le pari

Les points de terminaison conçus avec l'observabilité en tête dès le départ devraient être ceux qui se consolident correctement lors d'un refactoring — un graphe qui ne mentira pas sur les routes réellement porteuses de charge par rapport à celles qui ont simplement le plus de trafic historique. C'est la distinction qu'API Darwin est conçu pour révéler, avant qu'une vraie migration en paie le prix pour l'apprendre à la dure.
