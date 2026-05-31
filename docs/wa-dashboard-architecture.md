# WA Dashboard - Arsitektur & Dokumentasi

## Overview

WA Dashboard adalah platform SaaS standalone untuk mengelola operasional WhatsApp secara menyeluruh. Platform ini bersifat **multi-tenant** dan dapat digunakan oleh berbagai bisnis, termasuk Finku sebagai salah satu tenant.

WA Dashboard **tidak menyimpan data bisnis client**. Perannya adalah sebagai orchestrator antara user WhatsApp, AI (opsional per tenant), dan API masing-masing bisnis.

## Tanggung Jawab

| Kategori | Detail |
| --- | --- |
| **Broadcast** | Kirim pesan massal, promo, announcement |
| **CS Inbox** | Tim CS jawab chat masuk dari user |
| **Analytics** | Open rate, response time, delivery rate |
| **Multi-tenant** | Tiap bisnis punya workspace terpisah |
| **AI (per tenant)** | Chatbot, draft reply, summarize conversation |
| **WA Connection** | Manajemen koneksi WhatsApp Business API |

## Arsitektur Multi-Tenant

```text
WA Dashboard (SaaS)
|
+-- Tenant: Finku
|   +-- Broadcast (promo, notif massal)
|   +-- CS Inbox
|   +-- Analytics
|
+-- Tenant: Bisnis B (e-commerce)
|   +-- Broadcast
|   +-- CS Inbox
|   +-- AI Chatbot (cek order, produk)
|   +-- Analytics
|
+-- Tenant: Bisnis C
    +-- ...
```

Setiap tenant memiliki:

- Workspace terisolasi
- Koneksi WA Business API sendiri
- Konfigurasi AI sendiri (opsional)
- Data yang tidak bercampur dengan tenant lain

## Modul Utama

### 1. Broadcast Management

Fitur pengiriman pesan massal ke kontak WhatsApp.

**Fitur:**

- Buat dan jadwalkan broadcast
- Segmentasi kontak berdasarkan tag, label, atau filter
- Template message sesuai standar WA Business API
- Preview sebelum kirim
- Riwayat dan status broadcast: delivered, read, failed

**Flow:**

```text
Buat broadcast
      |
      v
Pilih template & segmen kontak
      |
      v
Jadwalkan / kirim langsung
      |
      v
WA Business API proses pengiriman
      |
      v
Monitor status di dashboard
```

### 2. CS Inbox

Antarmuka untuk tim CS menjawab pesan masuk dari user.

**Fitur:**

- Inbox terpusat semua chat masuk
- Assignment chat ke agent tertentu
- Label dan kategori conversation
- Quick reply / canned response
- Status: open, in progress, resolved
- AI-assisted draft reply (opsional)

**Role Management:**

| Role | Akses |
| --- | --- |
| Admin | Full access semua fitur |
| Supervisor | Lihat semua chat, assign agent |
| Agent | Hanya handle chat yang di-assign |

### 3. AI Integration (per Tenant)

Tenant non-Finku dapat mengaktifkan AI chatbot langsung dari WA Dashboard.

> **Catatan:** Finku tidak menggunakan fitur ini karena mengelola AI chatbot-nya sendiri di Admin Finku.

**Config AI per tenant:**

```json
{
  "tenant_id": "bisnis-b",
  "ai": {
    "enabled": true,
    "model": "claude-sonnet / gpt-4o",
    "system_prompt": "Kamu adalah asisten toko...",
    "handoff_to_cs": true,
    "handoff_trigger": "user request human / komplain"
  },
  "integrations": [
    {
      "name": "check_order",
      "description": "Cek status pesanan user",
      "endpoint": "GET https://api.bisnis-b.com/orders/{id}",
      "auth": {
        "type": "api_key",
        "header": "X-API-Key"
      }
    },
    {
      "name": "get_products",
      "description": "Lihat daftar produk",
      "endpoint": "GET https://api.bisnis-b.com/products",
      "auth": {
        "type": "api_key",
        "header": "X-API-Key"
      }
    }
  ]
}
```

**AI Handoff Flow:**

```text
User chat
    |
    v
AI merespons otomatis
    |
    +-- Pertanyaan bisa dijawab AI -> balas langsung
    |
    +-- Komplain / minta CS manusia
              |
              v
         Notif ke CS agent
              |
              v
         Agent ambil alih conversation
```

### 4. Analytics

Dashboard monitoring performa WhatsApp per tenant.

**Metrics:**

| Metric | Keterangan |
| --- | --- |
| Delivery rate | % pesan terkirim |
| Open rate | % pesan dibuka |
| Response time | Rata-rata waktu CS reply |
| Conversation volume | Jumlah chat per hari/minggu |
| Resolution rate | % chat berhasil diselesaikan |
| Broadcast performance | Per campaign |

### 5. Template Management

Kelola template pesan yang disetujui WhatsApp Business API.

**Fitur:**

- Buat dan submit template untuk approval Meta
- Monitor status approval
- Kategori template: marketing, utility, authentication
- Variable management (`{{1}}`, `{{2}}`, dan seterusnya)

## Tenant Config Structure

Setiap bisnis yang mendaftar WA Dashboard akan memiliki workspace dengan konfigurasi berikut:

```json
{
  "tenant_id": "unique-id",
  "business_name": "Nama Bisnis",
  "wa_config": {
    "phone_number_id": "...",
    "waba_id": "...",
    "access_token": "..."
  },
  "ai_enabled": false,
  "features": {
    "broadcast": true,
    "cs_inbox": true,
    "analytics": true,
    "ai_chatbot": false
  },
  "agents": [
    {
      "id": "agent-1",
      "name": "CS Team",
      "role": "agent"
    }
  ]
}
```

## Posisi Finku sebagai Tenant

Finku terdaftar di WA Dashboard, namun **hanya menggunakan fitur operasional**:

| Fitur | Finku Gunakan? | Keterangan |
| --- | --- | --- |
| Broadcast | Ya | Promo, notif massal ke user |
| CS Inbox | Ya | Tim CS Finku handle chat |
| Analytics | Ya | Monitor performa WA |
| AI Chatbot | Tidak | Dikelola di Admin Finku |
| Function calling | Tidak | Dikelola di Admin Finku |
| Auth user | Tidak | Dikelola di Admin Finku |

Pemisahan ini dilakukan karena AI chatbot Finku menyentuh **data finansial sensitif** yang lebih aman dikelola dalam ekosistem Finku sendiri.

## Arsitektur Teknis

```text
User WA
   |
   v
WhatsApp Business API
   |
   v
WA Gateway Service (WA Dashboard)
   |
   +---> Broadcast Engine
   |
   +---> CS Inbox Service
   |
   +---> AI Engine (untuk tenant yang aktifkan)
   |         |
   |         v
   |    Function Calling -> API Tenant
   |
   +---> Analytics Service
```

## Stack yang Direkomendasikan

| Komponen | Teknologi |
| --- | --- |
| Backend | Node.js / Go |
| Database | PostgreSQL (tenant data) |
| Queue | Redis / BullMQ (broadcast queue) |
| AI Engine | Claude / GPT dengan function calling |
| WA API | Meta Cloud API |
| Realtime inbox | WebSocket / SSE |
| Auth | JWT + multi-tenant middleware |
