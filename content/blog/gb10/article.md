---
title: Spark GB10 is local LLM's best bet
slug: gx-spark-gb10
description: >
  DGX Spark (GB10, 128 GB unified memory) is proving to be a capable platform for  self-hosting LLMs.
date: 2026-07-14
author: Sudhakar Daggubati
tags: 
  - vLLM
  - Codex
  - Claude
summary: >
 

featured: true
---

DGX Spark (GB10, 128 GB unified memory) is proving to be a capable platform for self-hosting LLMs.

When paired with a harness like Codex CLI or Claude Code, it supports multi-turn, multi-agent inference without a hard dependency on frontier models.

Most local LLM examples use Ollama or OpenCode, but I wanted to test the harnesses that OpenAI and Anthropic provide.

I gave the same task to both OpenAI Codex and Claude: build a RAG application following 12-factor app principles, with support for provider switching.

Both used a locally hosted qwen3-coder-30b-a3b-instruct-fp8 served via vLLM, with a huge 130k+ context window.

What worked

✓ Codex iteratively executed multi-turn tasks, processing close to 750k tokens over a few iterations.

✓ Both executed multi-turn, single-agent flows pretty decently.

✓ A frontier model review found many gaps, which the local model then worked through.

What didn't work

⊘ Claude didn't work when max token limits were set at 32k. It kept sending its system context, far exceeding the model limits.

⊘ Even clamping the context limits using a LiteLLM proxy didn't help. Claude Code appears designed around much larger, 200k-class context lengths.

⊘ Yet to get the Claude/OpenAI desktop apps working.

Takeaway

It looks like a 2-node GB10 cluster is pretty much enough for multi-turn, multi-agent inference.

One could offload deep reasoning and review to a dense reasoning model, local or remote, while utilizing open models cost-effectively, without frontier models sitting in the critical path.