# HARaMA Frontend - Quick Setup Guide

## What's Included

This is a complete Next.js 16 frontend for the HARaMA exam grading system with:
- ✅ Supabase JWT authentication (email + Google OAuth)
- ✅ Protected routes with middleware session refresh
- ✅ API client with automatic Bearer token injection
- ✅ Complete exam management workflow
- ✅ Submission upload (single & batch)
- ✅ AI grading review interface
- ✅ Analytics and results export

## Prerequisites

You need:
1. **Node.js 18+** installed
2. **Supabase account** (free tier at supabase.com)
3. **HARaMA backend** running (Go API on port 8080)

## Step 1: Supabase Setup

1. Go to [supabase.com](https://supabase.com)
2. Create a new project
3. Wait for initialization (5-10 minutes)
4. Copy these values from Project Settings → API:
   - `NEXT_PUBLIC_SUPABASE_URL`
   - `NEXT_PUBLIC_SUPABASE_ANON_KEY`

## Step 2: Environment Configuration

Create `.env.local` in the project root:

```env
NEXT_PUBLIC_SUPABASE_URL=https://your-project.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=your_anon_key_here
NEXT_PUBLIC_API_URL=http://localhost:8080
```

Replace:
- `your-project` with your Supabase project name
- `your_anon_key_here` with your actual anon key
- `http://localhost:8080` with your backend URL (if different)

## Step 3: Install & Run

```bash
# Install dependencies
npm install

# Start development server
npm run dev
```

The app will be at `http://localhost:3000`

## Step 4: Test Authentication

1. Visit `http://localhost:3000`
2. You'll be redirected to `/login`
3. Click "Sign up" to create a test account
4. Verify your email (check Supabase auth logs)
5. Log in with your credentials

## API Integration

The frontend automatically:
1. Gets JWT token from Supabase session
2. Adds `Authorization: Bearer <token>` to all API requests
3. Refreshes token via middleware on each request
4. Handles token expiration

All API calls go to `NEXT_PUBLIC_API_URL/api/v1/...`

## Available Pages

| Route | Purpose |
|-------|---------|
| `/login` | User login |
| `/signup` | User registration |
| `/dashboard` | Home / statistics |
| `/exams` | List all exams |
| `/exams/new` | Create exam wizard |
| `/exams/[id]` | Exam details |
| `/exams/[id]/upload` | Upload submissions |
| `/exams/[id]/submissions/[subId]` | Grade review |
| `/exams/[id]/results` | Results & export |
| `/exams/[id]/analytics` | Analytics charts |
| `/settings` | User settings |

## Common Issues

### "Auth Error" on login
- Check `.env.local` has correct Supabase credentials
- Verify email is confirmed in Supabase Auth
- Check browser console for exact error

### API calls return 401
- Verify `NEXT_PUBLIC_API_URL` is correct
- Ensure backend is running
- Check JWT token in DevTools Network tab

### Can't upload files
- Backend must be running
- Check file format (PDF or image only)
- Verify file size < 50MB
- Check browser console for errors

### Page shows "loading..."
- Backend might be slow
- Check network in DevTools
- Verify API URL is accessible

## Next Steps

1. **Backend Setup**: Make sure the Go backend is running at the configured URL
2. **Create Exam**: Click "New Exam" to create your first exam
3. **Upload Files**: Use the upload page to submit student answers
4. **Review Grades**: Check the grading review page for AI grades

## File Structure

```
/vercel/share/v0-project/
├── app/                    # Next.js app directory
│   ├── (auth)/            # Public auth pages
│   ├── (dashboard)/       # Protected dashboard
│   └── layout.tsx         # Root layout
├── components/
│   ├── layout/            # Sidebar, header
│   └── ui/                # shadcn components
├── lib/
│   ├── api.ts             # API client
│   └── supabase/          # Auth config
├── .env.local             # Configuration (create this)
└── package.json
```

## Authentication Details

### Login Flow
1. User enters email/password
2. Supabase authenticates and returns JWT
3. JWT stored in httpOnly cookie (secure)
4. Middleware refreshes token on each request
5. API calls include Authorization header

### OAuth (Google)
1. User clicks "Sign in with Google"
2. Redirected to Google consent screen
3. After approval, redirected to `/auth/callback`
4. Backend exchanges code for JWT session
5. User logged in

## Development Commands

```bash
# Start dev server
npm run dev

# Build for production
npm run build

# Start production server
npm start

# Type checking
npx tsc --noEmit

# Lint code
npm run lint
```

## Production Deployment

### On Vercel (Recommended)
1. Push to GitHub
2. Import repo in Vercel Dashboard
3. Add environment variables:
   - `NEXT_PUBLIC_SUPABASE_URL`
   - `NEXT_PUBLIC_SUPABASE_ANON_KEY`
   - `NEXT_PUBLIC_API_URL` (your backend URL)
4. Deploy

### On Other Platforms
Requires:
- Node.js 18+
- Environment variables set
- Backend API accessible

## Support Resources

- **API Docs**: See `API_DOCUMENTATION-2850w.md` in project
- **Supabase Docs**: https://supabase.com/docs
- **Next.js Docs**: https://nextjs.org/docs
- **shadcn/ui**: https://ui.shadcn.com

## Security Notes

- ✅ JWT tokens in httpOnly cookies (XSS safe)
- ✅ All API calls require authentication
- ✅ Middleware validates session
- ✅ No sensitive data in localStorage
- ✅ CORS configured via backend

## Performance Tips

1. The app uses React Server Components where possible
2. Images are optimized via Next.js Image component
3. Charts use Recharts for lightweight rendering
4. API calls are made only when needed

Ready to go! Start with `npm run dev` and navigate to `/login` 🚀
