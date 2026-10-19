# Takluz_TTS

<p align="center">
  <img src="resource/banner.png" alt="Takluz_TTS — voice on stream without Bits"/>
</p>

<p align="left">
  <img src="https://img.shields.io/badge/Go-1.23-00ADD8?logo=go&logoColor=white" alt="Go"/>
  <img src="https://img.shields.io/badge/Fiber-v2-00ADD8" alt="Fiber"/>
  <img src="https://img.shields.io/badge/MongoDB-Atlas-47A248?logo=mongodb&logoColor=white" alt="MongoDB"/>
  <img src="https://img.shields.io/badge/OBS-WebSocket-302EBA" alt="OBS"/>
  <img src="https://img.shields.io/badge/SAMMI-Core-orange" alt="SAMMI"/>
  <img src="https://img.shields.io/badge/license-MIT-blue" alt="License"/>
</p>

> Twitch channel-points TTS bot — gives high-loyalty viewers the same voice on stream as Bits donors, no payment required.

## What it is

Takluz_TTS turns Twitch channel-point redemptions and chat commands into spoken audio on stream — **Bits not required**. It bridges [SAMMI Core](https://sammi.solutions/) to five pluggable TTS engines (OpenAI, Google Cloud, Azure, Botnoi, Resemble AI) and plays the synthesized audio through OBS Studio's media controls, all from a Windows PC.

## Demo

![Demo clip 1](https://github.com/user-attachments/assets/a13dec69-70a5-4600-8870-987a531c1613)

## Architecture

```mermaid
flowchart LR
    subgraph local["Streamer PC (local)"]
        direction TB
        SAMMI["SAMMI Core<br/>Twitch ↔ HTTP bridge"]
        Bot["Takluz_TTS Bot<br/>Go + Fiber<br/>:4444"]
        OBS["OBS Studio<br/>WebSocket"]
    end

    subgraph cloud["Cloud Services"]
        direction TB
        TTS["TTS Engine<br/>OpenAI · Google<br/>Azure · Botnoi · Resemble"]
        DB[("MongoDB Atlas<br/>request log")]
    end

    Twitch["Twitch Chat<br/>& Channel Points"] --> SAMMI
    SAMMI -->|POST JSON<br/>userName, message| Bot
    Bot -->|HTTPS synthesize| TTS
    TTS -->|audio bytes| Bot
    Bot -->|play media| OBS
    OBS -->|stream| Twitch
    Bot -->|log| DB
```

## Why this exists

Off-the-shelf Twitch TTS bots gate voices behind Bits — paying for voice time. That ties participation to wallet, not engagement. A viewer who's been in chat for 200 days earns the same voice slot as someone who tipped once.

Takluz_TTS swaps Bits for **channel points** — Twitch-native currency every viewer already earns for free, just by being active. Now the high-loyalty viewer who never tips can speak on stream exactly like the big tipper. Streamer pays nothing extra; the platform's own engagement metrics decide who talks.

## Tech stack

- **Go 1.23** + Fiber v2 — HTTP server on `localhost:4444`
- **SAMMI Core** — Twitch ↔ HTTP bridge (handles channel points + chat commands)
- **OBS WebSocket** (`andreykaipov/goobs`) — plays synthesized audio into the stream
- **MongoDB Atlas** — logs every request for moderation
- **PyInstaller** — `setting.exe` config GUI for credentials and OBS password

## Quickstart

1. Install [SAMMI Core](https://sammi.solutions/) and [OBS Studio](https://obsproject.com/).
2. Enable OBS WebSocket server (Tools → WebSockets Server Settings → Enable).
3. Run `setting.exe`, paste your TTS provider API key and OBS WebSocket password.
4. Import the matching SAMMI deck JSON from `sammi/` (one per TTS engine).
5. Run `takluz_tts.exe`. The bot listens on `localhost:4444`.

## Supported TTS engines

| Engine | Strongest for | Notes |
|---|---|---|
| OpenAI | English, high quality | `tts-1` / `gpt-4o-mini-tts` |
| Google Cloud | Wide language coverage | Multi-voice |
| Botnoi | Thai pronunciation | Local speech team |
| Azure | Enterprise reliability | SSML support |
| Resemble AI | Custom cloned voices | Voice UUID required |

## Design choices worth flagging

- **Pluggable TTS** — adding an engine is one route handler (`/new_engine`) plus one config struct. The existing five were added without touching each other.
- **No vendor coupling** — OBS, Twitch, and every TTS provider sit behind small interfaces (`goobs`, Fiber routes, sound package). Each can be swapped.
- **Local-first** — the bot runs on the streamer's PC. Nothing is proxied through a third party.
- **Channel points over Bits** — chose the Twitch-native free path over the paid path because it lowers the bar for casual viewers and rewards loyalty.

## Contact

Discord: [takluz](https://discord.gg/whnxWt4NKj)
