---
title: "How Many Layers? Rethinking Medallion Architecture"
slug: rethinking-medallion-architecture
description: >
    Should Bronze, Silver, and Gold remain rigid physical boundaries, or should modern governance capabilities allow data architecture to follow real domain, ownership, lifecycle, compliance, and operational boundaries?
date: 2026-05-16
author: Sudhakar Daggubati
tags:
    - Databricks
    - Medallion Architecture
    - Data Architecture
    - Data Governance
    - Unity Catalog
summary: >
    Medallion architecture remains a useful model for representing data quality and lifecycle states, but its physical implementation need not become an
    architectural ritual. With mature governance capabilities such as ABAC, fine-grained access controls, and dynamic masking, separation can increasingly be driven by domain ownership, lifecycle, compliance, blast radius, SLAs, and operational requirements rather than rigid Bronze, Silver, and Gold layers.
featured: true
---

“𝐇𝐨𝐰 𝐦𝐚𝐧𝐲 𝐥𝐚𝐲𝐞𝐫𝐬?”

Does the **𝐌𝐞𝐝𝐚𝐥𝐥𝐢𝐨𝐧 𝐚𝐫𝐜𝐡𝐢𝐭𝐞𝐜𝐭𝐮𝐫𝐞** popularized by #Databricks (Bronze, Silver, Gold) still need to be applied with ritualistic rigidity, or should it become more flexible given 𝐦𝐚𝐭𝐮𝐫𝐞 𝐠𝐨𝐯𝐞𝐫𝐧𝐚𝐧𝐜𝐞 𝐜𝐚𝐩𝐚𝐛𝐢𝐥𝐢𝐭𝐢𝐞𝐬 such as ABAC, PBAC, row-level security, column-level security, and dynamic masking?

Medallion is a re-branded Data quality lifecycle model which is better executed with separate domains, workspaces, or accounts over a real boundary, ownership, lifecycle, sensitivity, capacity isolation, deployment process, or legal/compliance separation concerns.

### The question should not be “𝐡𝐨𝐰 𝐦𝐚𝐧𝐲 𝐥𝐚𝐲𝐞𝐫𝐬?” ⚠️

 It should be: 
 - what data needs to be preserved raw
 - who owns each state, who can access it
 - what breaks if it sits together? 
 - Can we reprocess from source easily?
 - Do different teams own raw, curated, and serving data?
 - Do we need separate deployment/security boundaries?

## Domain boundaries 

Whether a flat governance structure to mange that data quality progression suits better or a strict medallion/xyz segmentation should purely be on domain boundaries,domain ownership, SLA differences, blast-radius control, ingestion replay, data quality, quarantine, cost/capacity isolation, lifecycle/retention rules.

In a mature Databricks/UC setup, 𝐟𝐥𝐚𝐭 𝐥𝐨𝐠𝐢𝐜𝐚𝐥 𝐬𝐭𝐫𝐮𝐜𝐭𝐮𝐫𝐞 can be perfectly valid. The boundary can be enforced by identity + privileges + policies as guardrails.

Architectural/physical separation is useful when you do not fully trust policy consistency, identity hygiene, operational discipline, or blast-radius controls across all access paths.. and that's what for many decades architects worked on rather a single mold.

Had seen this rigid structure convulsing when a team/domain becomes 𝐛𝐨𝐭𝐡 𝐬𝐮𝐩𝐩𝐥𝐢𝐞𝐫 𝐚𝐧𝐝 𝐜𝐨𝐧𝐬𝐮𝐦𝐞𝐫, architectural purity starts losing to operational practicality and data mesh-medellian ritual flow became new friction within same team/domain.

Medallion still holds as a useful data-state pattern, but should its physical implementation be driven by real ownership, lifecycle, compliance, and operational boundaries 𝐫𝐚𝐭𝐡𝐞𝐫 𝐭𝐡𝐚𝐧 𝐚𝐫𝐜𝐡𝐢𝐭𝐞𝐜𝐭𝐮𝐫𝐚𝐥 𝐫𝐢𝐭𝐮𝐚𝐥? Is my understanding correct?