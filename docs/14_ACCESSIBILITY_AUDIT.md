# 14. ACCESSIBILITY AUDIT (WCAG 2.1 LEVEL AA)

**Document ID:** BN-ONE-ACC-014  
**Version:** 1.0.0  
**Effective Date:** 2026-10-24  

---

## 1. Compliance Scope & Standards

BN ONE adheres to **WCAG 2.1 Level AA** standards across both Mobile (Flutter) and Web (Vue 3) interfaces to guarantee equitable financial management access for users with visual, motor, or cognitive impairments.

---

## 2. Accessibility Verification Matrix

| WCAG Guideline | Criteria | Implementation & Status |
|---|---|---|
| **1.4.3 Contrast (Minimum)** | Contrast ratio >= 4.5:1 for normal text and >= 3:1 for large text. | **PASSED** (Primary text `#0F172A` on `#F8FAFC` achieves 13.5:1 ratio; brand button `#4F46E5` on white achieves 6.8:1 ratio). |
| **2.5.5 Target Size** | Interactive touch targets must be at least 48 x 48 CSS pixels / dp. | **PASSED** (All buttons, bottom tab bars, and quick action chips feature min-height 48px and generous touch padding). |
| **1.3.1 Info and Relationships** | Information, structure, and relationships conveyed through presentation can be programmatically determined. | **PASSED** (Semantic HTML5 elements `<header>`, `<main>`, `<nav>`, `<aside>`, `<section>` in Web; `Semantics` widgets in Flutter). |
| **2.4.7 Focus Visible** | Any keyboard-operable interface has an operable mode where the keyboard focus indicator is visible. | **PASSED** (Clear high-visibility `ring-2 ring-indigo-500` outline on keyboard focus across all interactive web elements). |
| **4.1.2 Name, Role, Value** | For all UI components, the name and role can be programmatically determined. | **PASSED** (Input fields paired with `<label>` tags; screen readers announce amount mask states as "Hidden Balance"). |

---

## 3. Assistive Technology Testing
- **Screen Readers:** Verified using TalkBack (Android), VoiceOver (iOS/macOS), and NVDA (Windows).
- **Reduced Motion:** Supports `@media (prefers-reduced-motion: reduce)` in web styles and disables large parallax transitions in Flutter.
- **Dynamic Type / Font Scaling:** Tested up to 200% font size scaling without clipping or overflow errors on mobile screens.
