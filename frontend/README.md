# HARaMA - AI-Powered Exam Grading Platform

A complete Next.js 16 frontend for the HARaMA (Handwritten Analysis and Mark Allocation) system. This application allows teachers to create exams, upload student answer sheets, and leverage AI-powered grading with manual review capabilities.

## Technology Stack

- **Framework**: Next.js 16 with App Router
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **UI Components**: shadcn/ui
- **Authentication**: Supabase Auth with JWT
- **API Integration**: RESTful API with JWT Authorization
- **State Management**: React hooks with fetch
- **Charts**: Recharts

## Features

### Authentication
- Email/password sign up and login
- Google OAuth integration
- JWT-based authorization for all API calls
- Automatic session refresh via middleware

### Exam Management
- Create exams with questions and rubrics
- Multi-step exam creation wizard
- Support for different answer types (short answer, essay, MCQ, diagram)
- Optional question groups (e.g., answer 11a OR 11b)
- Question rubric management

### Submission Upload
- Single and batch file upload
- Drag-and-drop support
- PDF and image file support
- Automatic student ID mapping
- Progress tracking

### Grading & Review
- View OCR-extracted student answers
- Review AI-generated grades with confidence scores
- Override grades with audit trail
- Navigate between submissions and questions
- View AI reasoning for grades

### Analytics & Results
- Grading status distribution
- Score trends over time
- Submission statistics
- Export results (CSV, Excel)
- Question-level insights

## Getting Started

### Prerequisites

1. **Supabase Project**: Create a project at [supabase.com](https://supabase.com)
2. **Environment Variables**: Copy and configure `.env.local`

### Installation

1. **Clone or extract the project**

2. **Install dependencies**:
```bash
npm install
```

3. **Configure environment variables**:
Create a `.env.local` file with:
```env
NEXT_PUBLIC_SUPABASE_URL=your_supabase_url
NEXT_PUBLIC_SUPABASE_ANON_KEY=your_supabase_anon_key
NEXT_PUBLIC_API_URL=http://localhost:8080
```

4. **Run the development server**:
```bash
npm run dev
```

Visit `http://localhost:3000` and you'll be redirected to login.

## Project Structure

```
app/
├── (auth)/              # Authentication pages
│   ├── login/
│   ├── signup/
│   ├── signup/success/
│   └── callback/        # OAuth callback
├── (dashboard)/         # Protected dashboard
│   ├── layout.tsx       # Dashboard layout with sidebar
│   ├── dashboard/       # Home dashboard
│   ├── exams/           # Exam management
│   ├── settings/        # User settings
│   └── [id]/            # Exam detail pages
├── page.tsx             # Root redirect
└── layout.tsx           # Root layout

components/
├── layout/              # Sidebar, header
├── ui/                  # shadcn components
└── ...                  # Feature components

lib/
├── api.ts               # API client with JWT auth
├── supabase/            # Supabase clients
│   ├── client.ts        # Browser client
│   ├── server.ts        # Server client
│   └── proxy.ts         # Session refresh
└── utils.ts             # Utilities

middleware.ts           # Session refresh middleware
```

## API Integration

The frontend connects to a Go backend API. Key endpoints:

### Exams
- `POST /api/v1/exams` - Create exam
- `GET /api/v1/exams` - List exams
- `GET /api/v1/exams/{id}` - Get exam details
- `POST /api/v1/exams/{id}/questions` - Add question

### Submissions
- `POST /api/v1/exams/{id}/submissions` - Upload single submission
- `POST /api/v1/exams/{id}/submissions/batch` - Batch upload
- `GET /api/v1/submissions/{id}` - Get submission status
- `POST /api/v1/submissions/{id}/trigger-grading` - Start grading

### Grading
- `GET /api/v1/submissions/{id}/grades` - Get grades
- `POST /api/v1/submissions/{id}/questions/{qid}/override` - Override grade
- `GET /api/v1/submissions/{id}/questions/{qid}/feedback` - Get feedback

### Analytics
- `GET /api/v1/analytics/grading-trends` - Grading statistics
- `POST /api/v1/exams/{id}/export` - Export grades

All requests include JWT token in Authorization header automatically.

## Authentication Flow

1. User signs up/logs in with email/password or Google OAuth
2. Supabase returns JWT access token
3. Middleware refreshes session on each page load
4. API client automatically includes JWT in all requests
5. Backend validates JWT and applies row-level security

## File Upload

### Single Upload
- Select one PDF or image file
- Specify student ID
- Automatically uploaded to backend

### Batch Upload
- Drag-drop or select multiple files
- Auto-map filenames to student IDs
- Manually adjust mappings if needed
- All files uploaded at once

Supported formats: PDF, JPG, PNG, JPEG

## Grading Workflow

1. **Upload**: Submit student answer sheets
2. **Process**: Backend OCR extracts text and images
3. **Grade**: AI evaluates answers using rubric
4. **Review**: Teacher reviews AI grades
5. **Override**: Manually adjust scores if needed
6. **Export**: Export final results

## Environment Variables

| Variable | Description |
|----------|-------------|
| `NEXT_PUBLIC_SUPABASE_URL` | Your Supabase project URL |
| `NEXT_PUBLIC_SUPABASE_ANON_KEY` | Supabase anonymous key |
| `NEXT_PUBLIC_API_URL` | Backend API base URL (default: http://localhost:8080) |

## Deployment

### Deploy to Vercel (Recommended)

1. Push code to GitHub
2. Import repository in [Vercel Dashboard](https://vercel.com)
3. Add environment variables in project settings
4. Deploy

### Deploy to Other Platforms

For Docker, Netlify, or other platforms:
1. Build: `npm run build`
2. Start: `npm run start`
3. Ensure Node 18+ is available
4. Set environment variables

## Features by Page

### /login & /signup
- Email/password authentication
- Google OAuth button
- Form validation
- Error handling

### /dashboard
- Welcome message
- Statistics cards
- Recent exams list
- Quick actions

### /exams
- List all exams
- Search and filter
- View/create actions

### /exams/new
- Multi-step wizard
- Basic info → Questions → Review
- Question management
- Real-time validation

### /exams/[id]
- Exam details
- Question accordion
- Rubric preview
- Action buttons

### /exams/[id]/upload
- File upload interface
- Single/batch modes
- Student ID mapping
- Progress tracking

### /exams/[id]/submissions/[subId]
- OCR result preview
- AI grade display
- Override interface
- Question navigation

### /exams/[id]/results
- Question breakdown
- Export options
- Score summary

### /exams/[id]/analytics
- Grading trends chart
- Status distribution pie chart
- Statistics cards
- Summary metrics

### /settings
- API configuration info
- Security & logout
- About section

## Best Practices

1. **Always use Supabase JWT**: Never hardcode API keys
2. **Handle errors gracefully**: Show user-friendly error messages
3. **Validate inputs**: Use form validation before submission
4. **Check authentication**: Middleware protects routes
5. **Cache data wisely**: Use React hooks for caching
6. **Test uploads**: Verify with small files first

## Troubleshooting

### "Not authenticated" error
- Check if user is logged in
- Verify JWT token in browser DevTools
- Ensure Supabase is configured correctly

### API calls failing
- Check `NEXT_PUBLIC_API_URL` is correct
- Verify backend is running
- Check network tab in DevTools
- Look for CORS errors

### Files not uploading
- Verify file format is supported (PDF/image)
- Check file size is under 50MB
- Ensure student IDs are filled in
- Check browser console for errors

### Grades not appearing
- Trigger grading manually in submission page
- Wait for backend processing
- Refresh page to see updates
- Check backend logs for processing errors

## Development

### Run development server
```bash
npm run dev
```

### Build for production
```bash
npm run build
```

### Run type checking
```bash
tsc --noEmit
```

## License

© 2024 HARaMA. All rights reserved.

## Support

For issues or questions:
1. Check the troubleshooting section
2. Review API documentation
3. Check backend logs
4. Contact support team
