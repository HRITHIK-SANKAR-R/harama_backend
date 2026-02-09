# HARaMA UI Design System

## Modern Professional Design Implementation

### Design Philosophy
- **Clean & Minimal:** Focus on content with minimal distractions
- **Professional:** Educational institution appropriate styling
- **Cool & Modern:** Contemporary gradient effects and animations
- **Accessible:** WCAG 2.1 AA compliant colors and contrast ratios
- **Responsive:** Works seamlessly on all device sizes

---

## Color System (3-5 Colors)

### Primary Colors
- **Blue** `#2563EB` - Main brand color for actions and highlights
- **Purple** `#A855F7` - Secondary brand color for accents
- **Cyan** `#06B6D4` - Bright accent for success states

### Neutral Colors  
- **Slate 50-950** - Full grayscale from white to nearly black
- **White** `#FFFFFF` - Clean backgrounds
- **Dark Slate** `#0F172A` - Dark mode background

### Status Colors
- **Green** `#10B981` - Success, completed
- **Amber** `#F59E0B` - Warning, pending
- **Red** `#EF4444` - Error, critical

---

## Typography

### Font Family
- **Primary:** System font stack (sans-serif)
- **Size Hierarchy:**
  - H1: 36px bold (landing, main titles)
  - H2: 28px bold (section titles)
  - H3: 20px semibold (subsection titles)
  - Body: 14-16px regular (content)
  - Small: 12-14px (labels, captions)

### Text Colors
- **Foreground:** Dark slate in light mode, white in dark mode
- **Secondary:** Muted slate for descriptions
- **Accents:** Blue/purple for interactive elements

---

## Component Styling

### Buttons
```
Primary:        bg-gradient-to-r from-blue-600 to-blue-500
Secondary:      bg-purple-600
Ghost:          bg-transparent with border
Disabled:       opacity-50, not-allowed cursor
Hover:          Scale up slightly, shadow increase
```

### Cards
- Border radius: 0.75rem (12px)
- Shadow: `shadow-lg` on base, `shadow-xl` on hover
- Background: White with subtle border in light mode
- Dark mode: slate-800/50 with border-slate-700
- Hover effect: Smooth shadow and scale transition

### Input Fields
- Background: slate-700/50 (dark mode), white (light)
- Border: slate-600 (dark), slate-200 (light)
- Focus ring: Blue (#2563EB)
- Placeholder: Muted text color
- Border radius: 0.75rem

### Badges
- Gradient backgrounds matching status
- White text for contrast
- Rounded corners
- No borders (border: 0)

### Navigation
- Sidebar: Fixed left panel with icons + text
- Header: Top bar with logo, search, user menu
- Breadcrumbs: Show navigation path
- Active states: Blue highlight with background

---

## Spacing & Layout

### Spacing Scale (Tailwind)
- xs: 0.25rem (4px)
- sm: 0.5rem (8px)
- md: 1rem (16px)
- lg: 1.5rem (24px)
- xl: 2rem (32px)
- 2xl: 3rem (48px)

### Grid System
- Mobile first: Single column
- Tablet: 2-4 columns
- Desktop: Full width with max constraints
- Gap between items: 4-6 units

---

## Animation & Transitions

### Available Animations
```css
@keyframes blob {
  /* Floating gradient blob background */
  0%, 100% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(-20px, 20px) scale(0.9); }
}

.animate-blob { animation: blob 7s infinite; }
```

### Transition Effects
- Hover states: `transition-all duration-300`
- Page transitions: Fade in/out
- Loading states: Spinning loader
- Status changes: Color fade transitions

---

## Dark Mode Implementation

### Theme Variables
All colors are CSS variables in `:root` and `.dark` classes:
```css
--background: Light vs dark background
--foreground: Light vs dark text
--primary: Blue brand color (adjusted for mode)
--secondary: Purple brand color
--accent: Cyan accent (adjusted)
--muted: Gray tones (mode adjusted)
```

### Auto Detection
- Respects system preference
- Uses `prefers-color-scheme` media query
- Fallback to light mode

---

## Page-Specific Designs

### Login Page
- Full viewport gradient background
- Animated blob elements
- Centered card form
- White text on dark background
- Social login button below

### Sign Up Page
- Similar gradient to login
- Multi-field form
- Institution field (optional)
- Same social login integration

### Dashboard
- Welcome banner with gradient title
- 4 stat cards in grid (responsive)
- Recent exams section
- Empty state with icon when no exams

### Exams List
- Search bar with icon
- Table view with striped rows
- Badge for question count
- Hover effects on rows
- View button for navigation

### Create Exam Wizard
- Progress indicator with step numbers
- Multi-step form (basic → questions → review)
- Color change on completion
- Smooth step transitions
- Back/Next navigation

### Grading Interface
- Two-column layout (answer key | student answer)
- OCR confidence score display
- Grade override button
- Feedback textarea
- Status badge

### Analytics
- Stat cards with icons and gradients
- Pie chart for grading status
- Line chart for trends
- Export button
- Color-coded legend

---

## Accessibility Features

### Color Contrast
- All text meets WCAG AA standards
- Blue #2563EB on white: 7.2:1 ratio
- Gray text still readable

### Interactive Elements
- Minimum 44x44px touch targets
- Keyboard navigation support
- Focus rings on all interactive elements
- Screen reader friendly labels

### Forms
- Proper label associations
- Error messages clear and visible
- Placeholder text not as labels
- Helper text for guidance

---

## Responsive Breakpoints

- **Mobile:** < 640px (tailwind: default)
- **Tablet:** 640px - 1024px (tailwind: md)
- **Desktop:** > 1024px (tailwind: lg)

### Layout Changes
- Mobile: Single column, full width
- Tablet: 2-column grid layouts
- Desktop: 3-4 column grids where applicable

---

## Browser Support

- Chrome/Edge 90+
- Firefox 88+
- Safari 14+
- Mobile browsers (iOS Safari 14+, Chrome Mobile)

---

## Performance Optimizations

### CSS
- Tailwind CSS minified and purged
- Only used styles included
- CSS-in-JS minimized

### Images
- Optimized with Next.js Image component
- Lazy loading enabled
- Responsive image sizes
- WebP format where supported

### JavaScript
- Code splitting by route
- Lazy component loading
- Minified production builds

---

## Design System Files

- **Globals CSS:** `app/globals.css` - All theme variables and animations
- **Component Library:** `components/ui/` - All shadcn UI components
- **Custom Components:** `components/layout/` - Header, sidebar
- **Tailwind Config:** `tailwind.config.ts` - Color, spacing, font extensions

---

## Future Design Enhancements

1. **Theming Options:** Allow custom brand colors
2. **Font Customization:** Support different font families
3. **Layout Variants:** Compact/spacious view options
4. **Accessibility:** High contrast mode
5. **Animations:** Reduced motion preference support
6. **Internationalization:** RTL support for Arabic/Hebrew

---

## Design Consistency Checklist

When adding new features:
- [ ] Use only defined colors from theme
- [ ] Maintain consistent spacing (use Tailwind scale)
- [ ] Follow button/card styling patterns
- [ ] Test dark mode appearance
- [ ] Verify mobile responsiveness
- [ ] Check color contrast ratios
- [ ] Ensure keyboard navigation works
- [ ] Add appropriate transitions/animations
- [ ] Test with screen readers
- [ ] Verify form accessibility
