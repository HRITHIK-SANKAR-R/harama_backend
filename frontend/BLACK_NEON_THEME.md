# HARaMA - Black Neon Theme Documentation

## Overview
HARaMA now features a premium black-themed UI with vibrant neon accents, smooth animations, and modern glass-morphism effects. This theme provides a professional yet cutting-edge aesthetic perfect for an AI-powered education platform.

## Color Palette

### Primary Colors
- **Neon Cyan**: `#00FFFF` (RGB: 0, 255, 255) - Primary accent and highlights
- **Neon Blue**: Used in gradients for primary elements
- **Neon Purple**: `#A855F7` (RGB: 168, 85, 247) - Secondary accent
- **Hot Pink**: `#FF0080` (RGB: 255, 0, 128) - Tertiary accent
- **Neon Green**: For success states
- **Neon Yellow**: For warnings

### Neutral Colors
- **Pure Black**: `#000000` - Background
- **Dark Gray**: `#1F2937` - Secondary backgrounds
- **Light Gray**: `#D1D5DB` - Text on dark backgrounds
- **White**: `#FFFFFF` - Primary text

## Design System

### CSS Variables (app/globals.css)
```css
--background: 0 0% 2%;        /* Pure black */
--foreground: 0 0% 98%;       /* Off white */
--primary: 180 100% 50%;      /* Neon Cyan */
--secondary: 280 100% 55%;    /* Neon Purple */
--accent: 320 100% 50%;       /* Hot Pink */
```

### Glassmorphism Effects
All cards and containers use:
- `backdrop-blur-md` or `backdrop-blur-2xl` - Frosted glass effect
- `bg-black/50` or `bg-black/80` - Semi-transparent black
- `border border-cyan-500/20` - Subtle neon borders
- `shadow-2xl shadow-cyan-500/20` - Neon glow shadows

## Advanced Animations

### 1. Blob Animation
```css
.animate-blob {
  animation: blob 7s infinite;
}
```
Floating, morphing blob shapes for background elements.

### 2. Neon Glow
```css
.animate-neon-glow {
  animation: neon-glow 2s ease-in-out infinite;
}
```
Pulsing glow effect perfect for glowing borders and icons.

### 3. Pulse Neon
```css
.animate-pulse-neon {
  animation: pulse-neon 2s ease-in-out infinite;
}
```
Smooth opacity and glow pulsing for loading states.

### 4. Scale In
```css
.animate-scale-in {
  animation: scale-in 0.5s ease-out;
}
```
Smooth entrance animation for cards and elements.

### 5. Shimmer Effect
```css
.animate-shimmer {
  animation: shimmer 2s infinite;
}
```
Light reflection sweep effect for premium feel.

### 6. Gradient Shift
```css
.animate-gradient-shift {
  animation: gradient-shift 3s ease infinite;
}
```
Animated gradient backgrounds that subtly shift colors.

### 7. Neon Border Glow
```css
.animate-neon-glow {
  box-shadow: 0 0 30px rgba(0, 255, 255, 0.8);
}
```
Intense glowing border effect for focus states.

## Component Styling Guidelines

### Buttons
```tsx
{/* Primary CTA Button */}
<Button className="bg-gradient-to-r from-cyan-500 to-blue-600 hover:from-cyan-400 hover:to-blue-500 text-black font-bold shadow-lg shadow-cyan-500/50">
  Action
</Button>
```

### Cards
```tsx
<Card className="border border-cyan-500/20 bg-black/50 backdrop-blur-md shadow-2xl hover:shadow-cyan-500/20 transition-all">
  {/* Content */}
</Card>
```

### Input Fields
```tsx
<Input className="bg-black/50 border-cyan-500/30 text-white focus:border-cyan-400 focus:ring-cyan-500/50 focus:ring-2" />
```

### Badges
```tsx
<Badge className="bg-gradient-to-r from-cyan-500 to-blue-600 text-black border-0 font-bold">
  Label
</Badge>
```

## Background Effects

### Grid Pattern
SVG overlay creates subtle grid pattern:
```tsx
<svg className="absolute inset-0 w-full h-full opacity-5">
  <defs>
    <pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse">
      <path d="M 40 0 L 0 0 0 40" fill="none" stroke="currentColor" strokeWidth="0.5" />
    </pattern>
  </defs>
  <rect width="100%" height="100%" fill="url(#grid)" className="text-cyan-500" />
</svg>
```

### Floating Blob Backgrounds
```tsx
<div className="absolute -top-40 -right-40 w-96 h-96 bg-cyan-500 rounded-full mix-blend-screen filter blur-3xl opacity-5 animate-blob" />
```

## Responsive Design
- Mobile-first approach
- All animations disabled on `prefers-reduced-motion`
- Touch-friendly button sizes (min 44px)
- Responsive grid: `grid-cols-1 md:grid-cols-2 lg:grid-cols-4`

## Dark Mode
The theme is always dark - there's no light mode toggle. This maintains consistent neon aesthetics.

## Performance Tips
1. Use `opacity` classes instead of `rgba()` for better performance
2. Limit animations on low-power devices with media queries
3. Use `will-change` sparingly for heavy animations
4. Optimize SVG backgrounds for file size

## Accessibility
- All text has sufficient contrast with black backgrounds
- Neon colors have WCAG AA contrast with black
- Animations respect `prefers-reduced-motion`
- Focus states clearly visible with cyan glow
- Screen readers ignore decorative blobs and grids

## File Changes
- `app/globals.css` - Theme variables and animations
- `app/(auth)/login/page.tsx` - Black neon login page
- `app/(auth)/signup/page.tsx` - Black neon signup page
- `app/(dashboard)/dashboard/page.tsx` - Neon dashboard
- `app/(dashboard)/exams/page.tsx` - Neon exams list
- All component styling uses new color variables

## Future Enhancements
- Add particle effects on hover
- Implement 3D transform effects on cards
- Add sound effects for interactions (optional)
- Create variant theme colors (purple, pink dominant)
- Add animation presets via CSS classes

---

**Last Updated:** February 7, 2026
**Theme Version:** 2.0 - Black Neon Edition
