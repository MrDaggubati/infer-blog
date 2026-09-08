---
title: From Consumer Chaos to Enterprise-Grade Network Segmentation
slug: network-segmentation
description: >
    Moving from a flat home network to a segmented, deny-by-default architecture
    using Proxmox VE, OPNsense, managed switching, VLANs, and isolated trust
    zones for compute, storage, IoT, media, and work infrastructure.
date: 2026-09-08
author: Sudhakar Daggubati
tags:
- Zero Trust
- OPNsense
- Proxmox
- Homelab
- Firewall
- Cybersecurity
summary: >
    A practical journey from an ISP router and flat network to a segmented infrastructure built around Proxmox VE, virtualized OPNsense, managed
    switching, and 802.1Q VLANs. The setup separates compute, storage, IoT, media, and work environments into explicit trust zones with deny-by-default
    policies, reducing lateral movement while bringing enterprise networking principles into a self-hosted environment.
featured: false
---

## 🛠️ 𝐅𝐫𝐨𝐦 𝐂𝐨𝐧𝐬𝐮𝐦𝐞𝐫 𝐂𝐡𝐚𝐨𝐬 𝐭𝐨 𝐄𝐧𝐭𝐞𝐫𝐩𝐫𝐢𝐬𝐞 𝐆𝐫𝐚𝐝𝐞 𝐍𝐞𝐭𝐰𝐨𝐫𝐤𝐢𝐧𝐠 segmentation

In 2026, an ISP router and flat network are a liability. 

I’ve overhauled my gateway with a Zero-Trust infrastructure built to withstand the era of autonomous bots and lateral AI threats

The Stack:

**𝐇𝐲𝐩𝐞𝐫𝐯𝐢𝐬𝐨𝐫:** Proxmox VE (2-node cluster with 120GB RAM pool).

**𝐅𝐢𝐫𝐞𝐰𝐚𝐥𝐥:** OPNsense running as a VM (Virtualized Networking).

**𝐒𝐰𝐢𝐭𝐜𝐡𝐢𝐧𝐠:** Netgear Managed Switch handling the 802.1Q VLAN trunking.

**𝐒𝐭𝐨𝐫𝐚𝐠𝐞:** Synology NAS integrated into a dedicated high-speed storage VLAN.

**𝐕𝐋𝐀𝐍 𝐒𝐞𝐠𝐦𝐞𝐧𝐭𝐚𝐭𝐢𝐨𝐧:** Moving from a flat network to isolated zones (LAN, IoT, Media, Work). Implementing "Deny by Default" means my compute infra can be truly isolated, chatty IoT gadgets in their echo chambers.

## 𝐓𝐡𝐞 𝐋𝐞𝐬𝐬𝐨𝐧

Building a network that is "VLAN-Aware" from the Proxmox bridge down to the managed switch is interesting. One incorrect "Access vs. Trunk" setting on a physical port can bring your entire network offline; the management VLAN is there to rescue you.

In an era of autonomous bots and AI-driven exploits, a flat network isn't just a risk—it's a target.

Thanks to the wonderful folks at Deciso (a Dutch firm) for forking it and managing a true OSS Version.

#Proxmox #OPNsense #Networking #CyberSecurity #SelfHosted #DevOps