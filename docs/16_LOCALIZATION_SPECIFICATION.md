# 16. LOCALIZATION & INTERNATIONALIZATION (i18n) SPECIFICATION

**Document ID:** BN-ONE-LOC-016  
**Version:** 1.0.0  
**Effective Date:** 2026-10-24  

---

## 1. Supported Languages & Locales

BN ONE v1.0 targets global personal finance users with Day-1 localization support across 10 major global and regional languages:

| Locale Code | Language Name | Native Name | Script Direction | Default Base Currency |
|---|---|---|---|---|
| `ko-KR` | Korean | 한국어 | LTR | KRW |
| `en-US` | English (US) | English | LTR | USD |
| `ja-JP` | Japanese | 日本語 | LTR | JPY |
| `zh-CN` | Chinese (Simplified) | 简体中文 | LTR | CNY |
| `zh-TW` | Chinese (Traditional) | 繁體中文 | LTR | TWD |
| `es-ES` | Spanish | Español | LTR | EUR |
| `fr-FR` | French | Français | LTR | EUR |
| `de-DE` | German | Deutsch | LTR | EUR |
| `vi-VN` | Vietnamese | Tiếng Việt | LTR | VND |
| `th-TH` | Thai | ไทย | LTR | THB |

---

## 2. Number & Date Formatting Architecture

1. **Intl Formatters:**
   - Numbers and dates formatted dynamically based on user's active locale setting.
   - For example, `1500000` KRW:
     - `ko-KR`: `₩1,500,000`
     - `en-US`: `₩1,500,000`
   - For USD `$1250.50`:
     - `en-US`: `$1,250.50`
     - `fr-FR`: `1 250,50 $` / `1 250,50 €`
2. **Date & Calendar Standards:**
   - Gregorian calendar standard across all platforms.
   - ISO-8601 strings (`2026-10-24T12:00:00Z`) used in API requests and serialized database storage.
   - Formatted for display using localized patterns (e.g., `YYYY년 MM월 DD일` in Korean, `MMM DD, YYYY` in English).

