# HARaMA - Deployment Guide

## Quick Start

### Prerequisites
- Node.js 18+ installed
- Supabase account with JWT auth configured
- Backend API running (see backend docs)

### Local Development

```bash
# Install dependencies
npm install

# Create .env.local with:
NEXT_PUBLIC_SUPABASE_URL=https://your-project.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=your_anon_key
NEXT_PUBLIC_API_URL=http://localhost:8080

# Run development server
npm run dev

# Open http://localhost:3000
```

---

## Deploy to Vercel

### Option 1: Using Vercel CLI

```bash
# Install Vercel CLI
npm i -g vercel

# Login to Vercel
vercel login

# Deploy
vercel

# Add environment variables when prompted:
# - NEXT_PUBLIC_SUPABASE_URL
# - NEXT_PUBLIC_SUPABASE_ANON_KEY
# - NEXT_PUBLIC_API_URL (production backend URL)
```

### Option 2: GitHub Integration

1. Push code to GitHub repository
2. Go to vercel.com and import repository
3. Add environment variables in project settings
4. Deploy button will be available

---

## Environment Variables

**Required:**
```
NEXT_PUBLIC_SUPABASE_URL        # Your Supabase project URL
NEXT_PUBLIC_SUPABASE_ANON_KEY   # Supabase anonymous key
NEXT_PUBLIC_API_URL             # Backend API URL
```

**Optional:**
```
NEXT_PUBLIC_SITE_URL            # For OAuth redirects (defaults to deployed URL)
```

---

## Configuration

### Supabase Setup

1. **Create Supabase project** at supabase.com
2. **Enable email authentication:**
   - Auth → Providers → Email
   - Enable Email/Password

3. **Enable Google OAuth:**
   - Auth → Providers → Google
   - Add Google OAuth credentials
   - Set redirect URL: `https://your-project.supabase.co/auth/v1/callback`

4. **Get credentials:**
   - Copy Project URL from Settings
   - Copy Anon Key from Settings → API

### Backend Configuration

1. **Set JWT token validation:**
   - Backend should verify tokens from `Authorization: Bearer <token>`
   - Token issued by Supabase JWT auth

2. **Add CORS configuration:**
   ```
   Allow origin: https://your-vercel-domain.vercel.app
   Allow credentials: true
   Allow headers: Authorization, Content-Type
   ```

3. **Test API endpoint:**
   ```bash
   curl -H "Authorization: Bearer <your-jwt-token>" \
        https://your-backend.com/api/v1/exams
   ```

---

## Security Checklist

- [ ] JWT tokens are validated on backend
- [ ] CORS headers configured correctly
- [ ] API accepts only authenticated requests
- [ ] Sensitive endpoints protected with role checks
- [ ] HTTPS enabled on all domains
- [ ] Environment variables marked as secret in Vercel
- [ ] OAuth redirect URLs whitelisted
- [ ] Rate limiting enabled on backend

---

## Performance Optimization

### Already Configured
- Image optimization with Next.js Image component
- CSS minification
- JavaScript code splitting
- Font optimization
- Static file caching

### Additional Optimizations
```bash
# Enable React Compiler (optional)
# In next.config.js:
# reactCompiler: true

# Build bundle analysis
npm run build -- --analyze
```

---

## Monitoring & Debugging

### Vercel Analytics
- Real-time analytics dashboard
- Performance metrics
- Error tracking
- Deployment history

### Logs
```bash
# View deployment logs
vercel logs [url]

# View function logs
vercel logs [url] --follow
```

### Browser DevTools
- Check Network tab for API responses
- Console for JavaScript errors
- Application tab for local storage/cookies

---

## Troubleshooting

### Issue: "Not authenticated" error

**Solution:**
1. Check Supabase credentials in `.env.local`
2. Verify JWT token is being sent in API headers
3. Check browser cookies for `sb-auth-token`

### Issue: API requests fail with 401

**Solution:**
1. Verify backend is expecting `Authorization: Bearer <token>`
2. Check token expiration in Supabase
3. Verify CORS headers from backend

### Issue: Google OAuth not working

**Solution:**
1. Verify Google OAuth credentials in Supabase
2. Check OAuth redirect URL matches deployment domain
3. Check browser console for OAuth error messages

### Issue: Slow page loads

**Solution:**
1. Check API response times
2. Verify backend is running efficiently
3. Check Vercel analytics for bottlenecks
4. Enable caching on frequently accessed endpoints

---

## Rollback Procedure

```bash
# View deployment history
vercel list deployments

# Rollback to previous deployment
vercel rollback [deployment-url]
```

---

## CI/CD Pipeline

### GitHub Actions (Optional)

```yaml
name: Deploy to Vercel

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: vercel/action@v4
        with:
          vercel-token: ${{ secrets.VERCEL_TOKEN }}
          vercel-org-id: ${{ secrets.VERCEL_ORG_ID }}
          vercel-project-id: ${{ secrets.VERCEL_PROJECT_ID }}
```

---

## Production Checklist

- [ ] Environment variables set in Vercel
- [ ] Backend API deployed and tested
- [ ] Supabase auth configured
- [ ] OAuth providers enabled
- [ ] CORS headers configured
- [ ] Error logging enabled
- [ ] Analytics configured
- [ ] Backup plan documented
- [ ] Team access granted
- [ ] Domain configured (if custom)

---

## Support

For issues:
1. Check IMPLEMENTATION_STATUS.md
2. Review Vercel logs
3. Check Supabase auth status
4. Test backend connectivity
5. Verify environment variables

---

## Next Steps

After deployment:
1. Share login credentials with test users
2. Conduct UAT (user acceptance testing)
3. Monitor error rates and performance
4. Collect user feedback
5. Plan feature updates
