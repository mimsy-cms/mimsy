# mimsy Landing Page Content & Structure

## Page Layout & Design Guidelines

### Visual Style
- **Aesthetic**: Clean, minimal, retro (early web/terminal inspired)
- **Typography**: Monospaced font for headers and code blocks
- **Color Scheme**: Simple, high contrast (black/white/minimal accent colors)
- **Graphics**: Simple ASCII art, basic diagrams, minimal icons

### Layout Structure
- Single-page design with clear sections
- Minimal navigation
- Focus on content hierarchy
- Mobile-responsive but desktop-first

---

## Section 1: Hero Section

### Main Headline
```
mimsy
A headless CMS for agencies who code
```

### Subheadline
```
Zero-downtime deployments. Code-first schema. TypeScript native.
Built for SvelteKit teams who want control over their CMS.
```

### Hero CTA
```
[View on GitHub] → https://github.com/[your-org]/mimsy
```

### Simple ASCII Art/Logo
```
    __________________
   /                  \
  |    mimsy v0.1.0    |
  |  [cms] [headless]  |
   \__________________/
```

---

## Section 2: The Problem (Agency Pain Points)

### Section Title
```
Why another CMS?
```

### Content
Every agency has been here:

**Payload CMS** → Over-engineered. Hooks everywhere. Setup nightmare.
**Strapi** → UI vs code conflicts. Sync hell. Inconsistent truth.
**PocketBase** → Migration chaos. Kubernetes nightmares.
**WordPress** → PHP baggage. UI modifications. Complex stack.

### Agency Reality Check
```
✗ Clients change requirements mid-project
✗ Zero-downtime deployments are non-negotiable  
✗ Multiple developers need consistent schemas
✗ TypeScript is your daily driver
✗ SvelteKit is your framework of choice
```

---

## Section 3: The Solution

### Section Title
```
How mimsy solves this
```

### Core Principles
```
1. CODE IS TRUTH
   Schema lives in TypeScript files
   No UI modifications allowed
   Git tracks all changes

2. ZERO DOWNTIME
   pgroll live migrations
   Old + new schemas coexist
   Rollback without service interruption

3. TYPESCRIPT NATIVE
   Type-safe resource definitions
   Compile-time validation
   SDK with full intellisense

4. SVELTEKIT FIRST
   Built by SvelteKit team for SvelteKit teams
   Native integration patterns
   Familiar development workflow
```

---

## Section 4: Technical Overview

### Architecture Diagram (Simple Text)
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Developer     │    │   mimsy CLI     │    │   Go Backend    │
│   (TypeScript)  │───▶│   (Schema)      │───▶│   (REST API)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                       │
                                                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Admin Panel   │    │   Your App      │    │   PostgreSQL    │
│   (SvelteKit)   │◀───│   (SvelteKit)   │◀───│   + pgroll      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Workflow Example
```javascript
// 1. Define your schema in TypeScript
// collections/blog.ts
export const BlogPost = {
  name: 'blog-posts',
  fields: {
    title: { type: 'string', required: true },
    slug: { type: 'string', unique: true },
    content: { type: 'richtext' },
    publishedAt: { type: 'date' }
  }
}

// 2. Process with CLI
$ mimsy build

// 3. Deploy with zero downtime
$ mimsy deploy
✓ Schema migration created
✓ Old API version: active
✓ New API version: active
✓ Traffic migration: complete
```

---

## Section 5: Key Features

### Feature Grid
```
┌─────────────────┬─────────────────┐
│ ZERO DOWNTIME   │ TYPE SAFETY     │
│ pgroll live     │ TypeScript      │
│ migrations      │ definitions     │
├─────────────────┼─────────────────┤
│ CODE FIRST      │ SVELTEKIT       │
│ Git-tracked     │ Native          │
│ schema          │ integration     │
└─────────────────┴─────────────────┘
```

### Feature Details

**Zero Downtime Deployments**
- Live migrations with pgroll
- Both old and new schemas active simultaneously
- Rollback without service interruption
- Perfect for agency production environments

**Code as Source of Truth**
- All schema changes happen in TypeScript files
- No UI-based modifications allowed
- Git tracks every change
- Perfect for team collaboration

**TypeScript Native**
- Type-safe resource definitions
- Compile-time validation
- Full IntelliSense support
- Catches errors before deployment

**SvelteKit Integration**
- Built by SvelteKit developers
- Native integration patterns
- Familiar development workflow
- Perfect for agency tech stacks

---

## Section 6: Quick Start

### Section Title
```
Get started in minutes
```

### Installation Steps
```bash
# 1. Install mimsy CLI
npm install -g mimsy-cli

# 2. Initialize project
mimsy init my-cms

# 3. Define your first collection
# collections/posts.ts
export const Posts = {
  name: 'posts',
  fields: {
    title: { type: 'string' },
    content: { type: 'richtext' }
  }
}

# 4. Build and deploy
mimsy build
mimsy deploy
```

### Main CTA
```
[View Full Documentation on GitHub] → https://github.com/[your-org]/mimsy
```

---

## Section 7: Status & Roadmap

### Current Status
```
STATUS: Early Development
VERSION: 0.1.0
STAGE: Architecture & Core Implementation
```

### Roadmap Preview
```
✓ Core architecture design
✓ Go backend foundation
⚡ TypeScript SDK development
⚡ SvelteKit admin panel
⚡ pgroll integration
◯ Beta release
◯ Production ready
```

---

## Section 8: Footer

### Links
```
GitHub: https://github.com/[your-org]/mimsy
Documentation: https://github.com/[your-org]/mimsy/docs
Issues: https://github.com/[your-org]/mimsy/issues
```

### Credits
```
Built with ❤️ by [your-team]
For agencies who code
```

---

## Technical Implementation Notes

### SvelteKit Component Structure
```
landing/src/routes/+page.svelte
├── HeroSection.svelte
├── ProblemSection.svelte
├── SolutionSection.svelte
├── TechnicalOverview.svelte
├── FeaturesSection.svelte
├── QuickStart.svelte
├── StatusSection.svelte
└── Footer.svelte
```

### Styling Guidelines
- Use monospaced fonts (Courier New, Monaco, or similar)
- Minimal color palette (black, white, one accent color)
- ASCII art and simple text diagrams
- Clean spacing and typography
- Terminal/retro aesthetic

### Content Delivery
- Single-page layout
- Progressive disclosure
- Clear section breaks
- Mobile-responsive but desktop-optimized
- Fast loading with minimal assets

---

## Content Tone & Voice

### Writing Style
- **Direct and honest**: Address real agency pain points
- **Technical but accessible**: Speak developer-to-developer
- **Confident but humble**: This is early-stage software
- **Problem-focused**: Lead with problems, follow with solutions

### Key Messaging
1. **Built by agencies, for agencies**
2. **Code-first approach eliminates common CMS headaches**
3. **Zero-downtime deployments are non-negotiable**
4. **TypeScript native for modern development workflows**
5. **SvelteKit integration because that's what agencies use**

This markdown file serves as the complete blueprint for the mimsy landing page, incorporating the retro aesthetic, agency-focused messaging, and technical depth that will resonate with the target audience.