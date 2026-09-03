---
title: MLflow as an AI Gateway 
slug: mlflow-ai-gateway-limits
description: >
    MLflow's AI Gateway offers solid tracing and basic budgets, but lacks dynamic routing, cache management, and cost-aware policy needed to govern coding agents like Claude Code and Codex at scale."
date: 2026-08-03
author: Sudhakar Daggubati
tags: 
    - MLflow
    - AI Gateway
    - LLMOps
summary: >
    MLflow's AI Gateway handles request tracing and basic budgets well, but it isn't yet built to govern coding agents at scale.Routing by capability, prompt, or intent, enforcing budgets across teams and applications, and managing cache and cost, that is a different layer of maturity, and one MLflow hasn't reached yet."

featured: false

---


Paired a local vLLM inference endpoint with MLflow's AI Gateway and quickly found it falls short when you try to govern coding agents like Claude Code and Codex at scale.

The moment you need to route by capability, prompt, intent, or tool and enforce budgets at the model, team, and application level simultaneously MLflow simply hasn't matured to that layer yet.

It gives you solid request tracing, basic budget caps, and simple guardrails. What it doesn't give you: cache management, cost-based routing, or provider-specific policy enforcement. 
Those aren't edge-case features for agent traffic, they're the core of what an AI gateway is supposed to do once usage moves past a handful of manual integrations.


##  Here's what's missing

- Dynamic routing across inference pools, instead of statically configured endpoints
- Rich traffic management , retries, failover, circuit breaking
- Semantic and exact-match cache management
- Fine-grained, provider and model-specific policies and restrictions
- Capability, prompt, intent, or tool-aware routing

One nuance worth separating out: budget enforcement and cost-based routing are not the same capability. A token-bucket budget stops a team from overspending. Cost-aware routing actively shifts a request to a cheaper or faster model based on live signals latency, load, price. MLflow does some of the first and none of the second. That second piece needs an actual routing policy engine, not just a proxy with counters attached.

If you're evaluating alternatives: 
- LiteLLM Proxy is the closest open-source comparator, with cost/latency-based routing and per-provider rate limits. 
- Portkey adds semantic caching and conditional routing by prompt metadata. 
- Kong and Envoy's AI gateway plugins lead with traffic management first, LLM-awareness second. 

Cloudflare AI Gateway is strong on caching and analytics, weaker on fine-grained policy.

### Cost-aware routing intelligence
```

        ^
   high |                    Portkey
        |                 o
        |
        |                        o LiteLLM
        |
        |
        |
        |
        |
        |
        |
        |           o Cloudflare
        |
        |
        |
        |                    o Kong / Envoy
        |
        |    o MLflow
    low +---------------------------------------------> 
        low                                      high
              Budget / quota enforcement
```


For now, most teams are running MLflow for tracing and eval, and a separate gateway for the control plane. Hopefully the standalone gateway catches up but until it does, that split seems to be where things land.

#MLflow #AIGateway #LLMOps #PlatformEngineering