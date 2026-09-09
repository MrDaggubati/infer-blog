---
title: The Self-Hosted Creed- Local LLMs
slug: local-llm-on-dev-box
description: >
 I have been a Self-hosting enthusiast for some time, exploring the stack that I want to learn. When the hype of LLMs hit the ceiling, I wanted to self-host that too, got a nice laptop with 8GB dGPU, wanted to use that at will, but got some issues
date: 2025-10-20
author: Sudhakar Daggubati
tags: 
  - vLLM
  - NVIDIA
summary: >
featured: false
---


#NVIDIA SPARK DGX Vs a homelab question would arise for many sooner or later.

I've been running local LLMs behind a Cloudflare Zero Trust tunnel with
GitHub OAuth (OIDC) authentication , all hosted on a type-1 hypervisor guest VM.

I wish to have a sleek system like DGX, but learning from issues that are cropping up in a capped system where expansion is not possible, I prefer to get a custom build; preferably a dual RTX 4090.

FOR now warming up a chunky laptop and exploring local LLMs and #ZerotrustSecurity.

Getting a stable context window of ~ 50 K - 64 K tokens on a type-1 hypervisor guest VM with 20GB RAM and 11 vCPUs where 8GB dGPU is offloaded from host to guest.

The context length of ~48-60k words / 6-7k lines looks small compared to the  Cloud-hosted LLM context window, yet it's manageable for small POC's and modularized code.

✅ What's works well
  - High-end server grade HWCONFIG 
  -  dGPU passthrough & ~17 logical cores
  - 100GB system RAM
  -  Zero trust security 
      https://chat.▓▓▓▓.sh
  -  Even small models provided  
      bigger therotical context.
  -  CPU offloading and hybrid training possibilities.
  - Datasets fits in memory, useful for CPU based inference. 
 -  Renting similar speced machine in cloud is quite expensive.

🚧⚠️ Things that are annoying. 

  - Fixed dGPU; GPU's physically soldered onto the motherboard, not upgradeable.
  - 𝐂𝐨𝐧𝐬𝐮𝐦𝐞𝐫-𝐠𝐫𝐚𝐝𝐞 𝐆𝐏𝐔𝐬 𝐚𝐫𝐞𝐧'𝐭 𝐬𝐮𝐩𝐩𝐨𝐫𝐭𝐞𝐝 𝐰𝐢𝐭𝐡 𝐯𝐆𝐏𝐔; 1 VM 1 GPU.
  - 4Bit quantised models plateaus at 4GB leaving remaining 4GB GPU idle
  - 𝐓𝐡𝐞𝐫𝐦𝐚𝐥 𝐧𝐨𝐢𝐬𝐞, 🔉Howls when large models exceed GPU specs.
 - Must reboot to claim GPU back to
   host.

💡Lessons learned:  
Could have got a dedicated dual-GPU RTX 4090 instead of a laptop.

For me, a dual RTX 4090 looks like a better option over SPARKX which doesn't expand when new hardware is available; 

Why not Apple!?, 
Most Android, Linux folks do not buy expensive 🍎😉.