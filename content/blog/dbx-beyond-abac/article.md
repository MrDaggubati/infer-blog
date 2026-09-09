---
title: Databricks ABAC governance wishlists
slug: databricks-abac-governance
description: >
    Databricks ABAC simplifies data governance through tags, inheritance,
    masking, and policy-based access. The next evolution could be
    topology-aware governance built around business domains, trust boundaries,
    security zones, and composable policy templates.
date: 2024-08-14
author: Sudhakar Daggubati
tags:
    - Databricks
    - Unity Catalog
    - ABAC
    - Data Governance
summary: >
    Databricks ABAC moves data governance from resource-by-resource access
    controls toward reusable, attribute-driven policies. Drawing parallels
    with the evolution of firewall security from individual rules to zones and identity-aware zero trust, this article explores a possible next step: topology-aware data governance. Business domains, trust boundaries,
    zone mappings, domain templates, and composable policies could provide an intent-driven governance layer capable of scaling across complex
    enterprises and the emerging agentic ecosystem.
featured: true
---

Databricks has reached an important milestone with ABAC governance. It provides a way to govern data assets bound to Databricks by combining long-established data domains, flexible tags, and ML-based classification into a practical ABAC capability.

ABAC answers an important question: **Should this user be allowed to see this data?**

But businesses typically think in terms of domains and relationships:

> **Should this domain interact with another domain — under which trust model, through which interfaces, and for which purposes?**

ABAC partially addresses these questions and is a strong first step toward better governance and policy enforcement. However, what domain-centric governance requires in addition is **governance orchestration**.

We had a similar problem with Kubernetes. Kubernetes provides containers, scheduling, networking, and other foundational primitives, but enterprises often need service meshes, operators, and topology orchestration to turn those primitives into a complete operational model.

Recently, I spent some time working with firewalls, network segmentation, and zone concepts. What I saw there got me thinking about how closely these ideas relate to ABAC — and how ABAC-based governance could eventually evolve toward **topology-based governance for the agentic era**.

## Background

I initially started with OPNsense and later moved to a Ubiquiti UDR7. Fortunately, I have only a few custom policies. At enterprise scale, however, the number of policies can easily reach thousands, making them increasingly difficult to understand, reason about, and manage.

Firewall appliance vendors addressed this complexity with sensible zone-based defaults. Users can define custom policies through a zone matrix rather than reasoning about every individual network endpoint.

The same concern applies here: zone-based policy models evolved considerably during the 2004–2010 period, and vendors such as Zscaler and Cisco later added another layer — identity- and context-based policies — as part of the broader zero-trust approach.

This is where I see a connection with Databricks' ABAC capabilities.

The domains are different, and so are the implementation details, but the underlying governance problem is surprisingly similar: **how do you define trust boundaries and policies without forcing users to manage every individual resource?**

## Data Segmentation

In network security, engineers work with packets, IP addresses, VLANs, VXLANs, ports, and network zones.

In data governance, data stewards and engineers work with datasets, domains, storage mounts, API endpoints, and gateway access.

Until recently, Unity Catalog primarily relied on roles and grants — concepts familiar to ETL, data warehouse, and OLTP practitioners. It has since been augmented with capabilities such as tagging and ML-based classification.

Recent Unity Catalog and ABAC changes help alleviate several operational burdens:

* table-by-table ACL sprawl
* view proliferation
* manual masking
* duplicated governance logic

The key architectural shift enabled by ABAC is:

> **Govern once, apply everywhere.**

Teams can also build similar custom ABAC layers using tools such as Casbin, with a proxy between consumers and data sources enforcing ABAC policies.

## The Next Evolution

In my view, Databricks and its users could evolve toward **enterprise-scale, domain- and zone-based governance**, with zone maps and templates similar in spirit to approaches used by vendors such as Zscaler.

The goal would be to simplify policy management while supporting increasingly complex governance requirements.

Databricks already provides capabilities such as tags, inheritance, masking, row filters, and time-based access controls. The next step, in my view, is to move beyond resource-centric policy management toward a model that understands **organizational topology and business intent**.

The agentic era will increasingly require:

* topology-based domain trust maps
* default-deny governance planes
* inheritance and precedence rules
* composable policies
* explicit trust boundaries
* policies expressed in terms of business purpose and domain relationships

With today's ABAC model, you still tend to think in terms of **tables, columns, tags, and an object explorer**.

What is missing is an interactive governance topology where teams can reason about:

* security domains
* trust boundaries
* governance intent
* domain-to-domain relationships
* permitted purposes and interfaces

## Topology-Aware Governance

A topology-aware governance model could map business domains into discrete **business zones**, with zone mappings and domain templates that simplify both policy expression and enforcement.

Databricks with ABAC is already close to this direction. The next step is less about adding another policy primitive and more about adding an orchestration layer around the existing ones:

* **topology-oriented UX**
* **zone-oriented policy orchestration**
* **domain templates**
* **policy composition**
* **inheritance and precedence**

Many data engineers do not have deep domain expertise, while business users often do. Both groups therefore need an intuitive UX that allows them to collaboratively define zone-oriented domain mappings and templates.

Those mappings can then feed the policy engine, which generates inheritable and customizable policies.

The architectural transition would be:

> **From policies attached directly to resources → policies generated from organizational topology.**

Developers should be able to define templates or zone mappings and drag and drop them into security domains. Default policies would then apply automatically, while teams could add custom policies based on the specific context of a domain.

This is not intended to represent any particular organization's actual domain topology. It is a conceptual mapping model.

## Intent / Topology-Based Governance

The conceptual progression is:

## *Intent /Topology based* 
![Intent /Topology based](images/beyond-abac-dbx.png)


**Catalog → Schema → Table → Domain → Purpose → Trust Boundary**

Or, architecturally:

**Resource-attached policies → Policies generated from organizational topology**

An intent- and topology-driven governance model seems increasingly necessary.

There are already niche products that address parts of this scope. Some are mature, while others are evolving on top of existing platforms. Some solutions outperform Unity Catalog in specific areas, but no single platform fully addresses the entire problem space yet.

Personally, I would prefer a third-party or CNCF-oriented solution that acts as a **governance proxy and orchestration layer**, providing cross-domain and potentially cross-organizational governance rather than embedding the entire model into every individual vendor platform.

The opportunity is to move governance from **managing resources** to **expressing organizational intent**.

ABAC is an important step in that direction. The next step may be **topology-aware, intent-driven governance** — where policies are not merely attached to data assets, but derived from how domains relate to one another, what they are trusted to do, and for which purposes.
