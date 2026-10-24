# 03. UX/UI DESIGN SPECIFICATION & DESIGN TOKENS

**Document ID:** BN-ONE-UX-003  
**Version:** 1.0.0  
**Effective Date:** 2026-10-24  

---

## 1. Design Tokens & Color Palette

BN ONE incorporates a high-trust, modern financial aesthetic built upon an indigo/slate palette. The interface balances high data density with generous white space and legible typography.

### 1.1 Color Tokens

| Token Name | Hex Code | Flutter / CSS Variable | Purpose |
|---|---|---|---|
| **Brand Primary** | `#4F46E5` | `AppColors.brandPrimary` / `--color-primary` | Primary buttons, brand headers, active tab icons |
| **Brand Accent** | `#06B6D4` | `AppColors.brandAccent` / `--color-accent` | Secondary highlights, badges, charts |
| **Success / Inflow** | `#10B981` | `AppColors.success` / `--color-success` | Income entries, positive net worth deltas, budget surplus |
| **Danger / Outflow** | `#EF4444` | `AppColors.danger` / `--color-danger` | Expense entries, budget overruns, deletion actions |
| **Warning** | `#F59E0B` | `AppColors.warning` / `--color-warning` | Disclaimer badges, subscription renewal alerts |
| **Surface Dark (Bg)** | `#0B0F19` | `AppColors.backgroundDark` / `--bg-dark` | Deep dark mode background |
| **Surface Light (Bg)**| `#F8FAFC` | `AppColors.backgroundLight` / `--bg-light` | Default light mode canvas |
| **Card / Container** | `#FFFFFF` / `#111827` | `AppColors.cardSurface` | High-elevation container cards |
| **Text Primary** | `#0F172A` / `#F8FAFC` | `AppColors.textPrimary` | Main titles, amounts, ledger figures |
| **Text Secondary** | `#64748B` / `#94A3B8` | `AppColors.textSecondary` | Subtitles, timestamps, disclaimers, captions |

---

## 2. Typography & Hierarchy

- **Primary Font Family:** Pre-system / Inter, `-apple-system`, `BlinkMacSystemFont`, `Segoe UI`, `Roboto`, `Helvetica Neue`, sans-serif.
- **Monospace / Numeric Font:** `ui-monospace`, `SFMono-Regular`, `Menlo`, `Monaco`, `Consolas`, monospace (for currency tables and ledger accounts).

| Hierarchy | Size (Web) | Size (Flutter) | Weight | Line Height | Usage |
|---|---|---|---|---|---|
| **Display 1** | 36px | 32sp | 900 (Black) | 1.15 | Net Worth Hero figure on Dashboard |
| **Heading 1** | 24px | 22sp | 700 (Bold) | 1.25 | Page titles, major category headers |
| **Heading 2** | 18px | 18sp | 600 (Semi-bold) | 1.35 | Section cards, modal headers |
| **Body Large** | 16px | 15sp | 500 (Medium) | 1.50 | Primary form fields, transaction items |
| **Body Regular** | 14px | 14sp | 400 (Regular) | 1.50 | Standard descriptions, secondary metrics |
| **Caption / Legal** | 12px | 11sp | 500 (Medium) | 1.40 | Legal disclaimers, timestamps, badges |

---

## 3. Accessibility & Touch Targets

1. **Touch Target Size:** All touchable interactive elements (buttons, segmented tabs, chips, form controls) have an explicit minimum hit boundary of **48 x 48 dp/px**.
2. **Contrast Ratios:**
   - Text Primary vs Background: **13.5:1** (Exceeds WCAG AAA).
   - Brand Primary vs Background: **6.8:1** (Exceeds WCAG AA).
   - Text Secondary vs Background: **4.6:1** (Exceeds WCAG AA minimum 4.5:1).
3. **Screen Reader Semantic Semantics:**
   - All amounts masked with privacy toggles read as `"Hidden Balance"` to screen readers.
   - All disclaimers feature `role="note"` and `aria-label="Legal Notice"`.

---

## 4. Privacy & Masking Feature UX

- Users can toggle the **Privacy Eye Icon** (located on the Net Worth hero card or app header).
- When toggled, all sensitive monetary balances immediately mask to `₩••••••`, `$••••••`, or `€••••••`.
- This ensures users can open the app in public transit, coffee shops, or airports without exposing their liquid wealth to shoulder-surfing.

