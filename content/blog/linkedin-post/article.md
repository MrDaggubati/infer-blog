---
title: Automted blog articles from a .md file
slug: markdown-posts
description: >
 Sometimes the simplest runtime is the one you don't need to run.

Static-first, not static-only: deploy only the runtime you actually need.
date: 2026-09-08
author: Sudhakar Daggubati
tags: 
  - go
  - automated
summary: >
  A static-first web portal built with Go and HTMX, where content, templates, and Markdown are compiled into plain HTML, CSS, JavaScript, and assets eliminating the need for a continuously running application server. The same codebase can later support dynamic, server-side features when needed, without abandoning the static deployment model. The result is a portal that remains simple by default while allowing capabilities like WebRTC, chatbots, or other services to evolve independently with their own runtimes.

featured: true
---

Websites often become application platforms too early. Others stay static for too long.

Add a 𝗯𝗹𝗼𝗴 and suddenly you need a 𝗰𝗼𝗻𝘁𝗲𝗻𝘁 𝘀𝘆𝘀𝘁𝗲𝗺.
Add 𝘀𝗲𝗮𝗿𝗰𝗵 and you need 𝗔𝗣𝗜𝘀.

Then deployment brings containers, runtime config health checks, monitoring and another service that needs to stay alive.

Wanted to see how much of that runtime we could simply avoid while rebuilding infer origins portal; was long overdue primarily due to lack of any front-end skills; definitely not fan of that complex world of UX :-).

I wanted my company portal  to sit somewhere in between: static by default, but designed to become an application when a capable runtime attached.

### what it comprises

- ✅ a 𝗚𝗼 + 𝗛𝗧𝗠𝗫  based web portal, still static-first web portal butbuilt with Go and HTMX, where 
     content, templates, and Markdown are compiled into plain HTML, CSS, JavaScript, and assets eliminating the need for a continuously running application server. 
- 🔄 The same codebase can later support dynamic, server-side features when needed, without abandoning 
     the static deployment model.
- ⏰ A portal that remains simple by default while allowing capabilities like WebRTC, 
     chatbots, or other services to evolve independently with their own runtimes.
 

> What gets deployed doesn't need a Go server running behind it.

But we're not locking ourselves into a static site either. If we need server-side behavior later, the same code base can run as 𝗚𝗼 + 𝗛𝗧𝗠𝗫 without throwing away the content or rendering model.

## Static-First, Not Static-Only

- **Adding WebRTC?** Make it another page with its own deployment model while the main portal remains static.
- **Adding a chatbot?** Use the same approach: another capability with its own deployment target.

There is no reason for the entire portal to inherit the runtime requirements of every feature.

**Sometimes the simplest runtime is the one you don't need to run**


No reason the whole portal needs to inherit its runtime."Sometimes the simplest runtime is the one you don't need to run".

Let me know what you think. The UX, copy, and CSS work is hell and not  my zone 😉


#Go #HTMX #PlatformEngineering #WebArchitecture


