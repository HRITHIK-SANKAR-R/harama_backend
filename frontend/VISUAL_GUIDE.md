# HARaMA Visual Design Guide

## 🎨 Design Philosophy
Black Neon Modern - Premium, Professional, Cutting-Edge

## 🎯 Color Palette

### Primary Neon Colors
| Color | Hex | HSL | Usage |
|-------|-----|-----|-------|
| Neon Cyan | #00FFFF | 180° 100% 50% | Primary actions, borders, highlights |
| Neon Purple | #8855FF | 260° 100% 55% | Secondary accents, gradients |
| Hot Pink | #FF0080 | 320° 100% 50% | Tertiary accent, hover states |
| Neon Green | #00FF00 | 120° 100% 50% | Success states |
| Neon Yellow | #FFFF00 | 60° 100% 50% | Warnings |
| Neon Blue | #0080FF | 220° 100% 50% | Gradients, secondary |

### Neutral Colors
| Color | Hex | HSL | Usage |
|-------|-----|-----|-------|
| Pure Black | #000000 | 0° 0% 0% | Background |
| Dark Gray | #1F2937 | 0° 0% 18% | Cards, containers |
| Medium Gray | #4B5563 | 0° 0% 25% | Secondary text |
| Light Gray | #D1D5DB | 0° 0% 82% | Primary text |
| Off White | #F5F7FA | 0° 0% 98% | Headings |

## 🎬 Animation Library

### 1. **Blob** (7s, Infinite)
Floating, morphing background shapes
```
Location: Background layers
Opacity: 5-20%
Blend Mode: screen/multiply
```

### 2. **Neon Glow** (2s, Infinite)
Pulsing border/icon glow
```
Intensity: 10px → 30px
Color: Cyan/Purple/Pink
Easing: ease-in-out
```

### 3. **Pulse Neon** (2s, Infinite)
Loading state pulsing
```
Effect: Opacity + glow
Intensity: 0.8 → 1
Used For: Loaders, icons
```

### 4. **Scale In** (0.5s, Ease-Out)
Element entrance zoom
```
From: scale(0.95), opacity(0)
To: scale(1), opacity(1)
Stagger: 50-100ms delays
```

### 5. **Shimmer** (2s, Infinite)
Premium light reflection sweep
```
Direction: Left to right
Gradient: rgba(255,255,255, 0-0.2-0)
Size: 200% width
```

### 6. **Gradient Shift** (3s, Infinite)
Animated background gradient
```
Movement: 0% → 100%
Colors: Rotate through palette
Position: 50% to 50%
```

### 7. **Float Up** (3s, Infinite)
Floating particle effect
```
Movement: 0px → -20px
Opacity: 0 → 1 → 0
Easing: ease-in
```

### 8. **Slide In** (0.5s, Ease-Out)
Directional entrance
```
From: translateX(100px)
To: translateX(0)
Opacity: 0 → 1
```

## 🧱 Component Styling

### Buttons
**Primary CTA**
```css
Background: linear-gradient(to right, #00FFFF, #0080FF)
Text: Black (font-bold)
Shadow: 0 0 20px rgba(0, 255, 255, 0.5)
Hover: Lighter gradient, increased shadow
Radius: 0.5rem
Height: 2.75rem (h-11)
```

**Secondary**
```css
Background: rgba(0, 0, 0, 0.5)
Border: 1px solid rgba(168, 85, 247, 0.5)
Text: White
Hover: Background more opaque
```

### Cards
```css
Background: rgba(0, 0, 0, 0.5)
Border: 1px solid rgba(0, 255, 255, 0.2)
Backdrop: blur(12px)
Shadow: 0 25px 50px rgba(0, 255, 255, 0.2)
Radius: 0.75rem
```

### Input Fields
```css
Background: rgba(0, 0, 0, 0.5)
Border: 1px solid rgba(0, 255, 255, 0.3)
Text: White
Focus: 
  - Border: #00FFFF
  - Ring: 2px solid rgba(0, 255, 255, 0.5)
```

### Badges
```css
Background: linear-gradient(to right, #00FFFF, #0080FF)
Text: Black (font-bold)
Border: none
Radius: 0.375rem
```

## 🎞️ Page Layouts

### Login/Signup Pages
```
┌─────────────────────────────────────────┐
│  [Animated Grid Background]             │
│  [3 Floating Blobs]                     │
│                                         │
│      ┌──────────────────────────┐      │
│      │  Logo + Brand Text       │      │
│      │  (Cyan-Purple Gradient)  │      │
│      │                          │      │
│      │  [Neon Glowing Card]     │      │
│      │  ├─ Email Input          │      │
│      │  ├─ Password Input       │      │
│      │  ├─ Sign In Button       │      │
│      │  ├─ Divider              │      │
│      │  ├─ Google Button        │      │
│      │  └─ Sign Up Link         │      │
│      │                          │      │
│      └──────────────────────────┘      │
│                                         │
└─────────────────────────────────────────┘
```

### Dashboard Page
```
┌─────────────────────────────────────────┐
│  [Grid Background + Blobs]              │
│                                         │
│  Dashboard (Cyan-Purple Gradient)       │
│  Welcome back, [User]                   │
│                                         │
│  ┌────────────┐ ┌────────────┐         │
│  │ Total Exams│ │Pending Rev.│         │
│  │    42      │ │     7      │         │
│  └────────────┘ └────────────┘         │
│  ┌────────────┐ ┌────────────┐         │
│  │ Auto Graded│ │Avg Confid. │         │
│  │     35     │ │   95.2%    │         │
│  └────────────┘ └────────────┘         │
│                                         │
│  Recent Exams          [New Exam ➜]    │
│  ┌──────────────────────────────────┐  │
│  │ Math Midterm  • 30 Q    12/25/26 │  │
│  └──────────────────────────────────┘  │
│                                         │
└─────────────────────────────────────────┘
```

### Exams List Page
```
┌─────────────────────────────────────────┐
│  [Grid Background]                      │
│                                         │
│  Exams                    [New Exam]    │
│  [Search Box]                           │
│                                         │
│  ┌───────────────────────────────────┐  │
│  │ Title │ Subject │ Q │ Date │ View │  │
│  ├───────────────────────────────────┤  │
│  │ Exam1 │ Math    │30│12/25│ ➜    │  │
│  │ Exam2 │ Science │25│12/24│ ➜    │  │
│  │ Exam3 │ English │40│12/23│ ➜    │  │
│  └───────────────────────────────────┘  │
│                                         │
└─────────────────────────────────────────┘
```

## 📱 Responsive Behavior

### Mobile (320-640px)
- Single column layout
- Full-width cards
- Stacked stat cards
- Touch-friendly buttons (44px min)
- Reduced animation complexity

### Tablet (641-1024px)
- 2-column grids
- Side-by-side cards
- Full animations enabled
- Optimized spacing

### Desktop (1025px+)
- 4-column grids
- Maximum visual complexity
- Advanced animations
- Perfect spacing and alignment

## 🌈 Interactive States

### Hover States
```css
Card: 
  - Border: Fully opaque cyan
  - Shadow: Increased cyan glow
  - Transform: Slight scale increase (optional)

Button:
  - Background: Lighter gradient
  - Shadow: Larger, more intense
  - Transform: Scale 1.02

Input:
  - Border: Cyan highlight
  - Shadow: Cyan glow
  - Ring: Visible focus indicator
```

### Focus States
```css
All Interactive:
  - Ring: 2px cyan
  - Outline: None (handled by ring)
  - Box-shadow: Cyan glow
```

### Active States
```css
Buttons:
  - Opacity: 0.9
  - Transform: Scale 0.98
  - Duration: 50ms
```

## ✨ Special Effects

### 1. Glassmorphism
Semi-transparent frosted glass effect with blur
```css
background: rgba(0, 0, 0, 0.5);
backdrop-filter: blur(12px);
border: 1px solid rgba(0, 255, 255, 0.2);
```

### 2. Neon Glow Shadows
```css
box-shadow: 0 0 20px rgba(0, 255, 255, 0.5),
            inset 0 0 10px rgba(0, 255, 255, 0.1);
```

### 3. Gradient Text
```css
background: linear-gradient(to right, #00FFFF, #8855FF, #FF0080);
-webkit-background-clip: text;
-webkit-text-fill-color: transparent;
background-clip: text;
```

## 🎯 Visual Hierarchy

### Text Sizes & Weights
```
H1 (Headings): 48px, font-black, gradient
H2 (Section): 32px, font-bold, white
H3 (Card Title): 18px, font-bold, white
H4 (Label): 14px, font-semibold, white
Body: 14-16px, font-normal, light gray
Small: 12px, font-normal, medium gray
```

### Color Intensity
```
Background: 0% (pure black)
Secondary BG: 5% (very dark)
Cards: 50% opacity (semi-transparent)
Borders: 20% opacity (subtle)
Text: 98% opacity (bright)
Accents: 100% opacity (pure neon)
```

## 🚀 Performance Notes

1. **GPU Acceleration** - All animations use `transform` and `opacity`
2. **Will-change** - Used strategically on animated elements
3. **Backdrop Filter** - Hardware accelerated in modern browsers
4. **SVG Grid** - Optimized for minimal file size
5. **Lazy Loading** - Animations trigger on interaction

## 📊 Design Metrics

| Metric | Value |
|--------|-------|
| Border Radius | 0.75rem (12px) |
| Padding | 1.5rem (24px) |
| Margin | 1rem-2rem (16-32px) |
| Gap/Spacing | 1rem (16px) |
| Shadow Blur | 12-30px |
| Animation Duration | 0.5s-7s |
| Transition Duration | 200-500ms |
| Opacity Ranges | 5%-100% |

---

**Last Updated:** February 7, 2026
**Version:** 2.0 Black Neon Edition
**Status:** Production Ready
