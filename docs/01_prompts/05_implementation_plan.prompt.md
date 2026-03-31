# Prompt: Generate Phased Implementation Plan

> **How to use:** Complete Prompts 01–04 first. Fill in every `[[PLACEHOLDER]]`, then paste this entire file into your AI tool. Save the output to `docs/02_outputs/05_implementation_plan.md`.

---

## Instructions to AI

You are a senior engineering lead. Generate a **phased implementation plan** that takes the project from zero to a production-ready v1. The plan must be granular enough that an AI coding agent can execute each task independently.

---

## Project Input

**Project name:** e-commerce-site
**Backend:** Go 1.26 / Gin Web Framework
**Frontend:** NextJS 16 Static Site Generation
**Database:** Postgres
**Auth strategy:** JWT
**Deployment target:** AWS
**CI/CD:** None

**v1 Features:**
1. Landing page for the website. The landing page will have a horizontal carousel. After the carousel there will be product cards
2. Product list page with the following requirements:
    - Each product will have the following fields categories, Level, theme, and size apart from the usual discount, and stock fields.
    - Users can filter products based on categories, Level, theme, and size. There should a filter section on the left side of the page
    - The current search criteria should be maintained in the URL as query params
    - Users can add product to his cart
    - These are the list of categories: Colouring Books, Drawing Books, Journals, Calendars, Sketch Pads, Planners, Stained Art, Fabrii Art, Note Book, Sticker Book, Sticker Colouring Book, and Activity Cards.
3. Navbar requirements:
    - There should be a search bar at the top right, where users can search for products
    - Here are the list of links that the navbar will have: Book Store, Wall of fame, Contact us, Printing service, and Printing supplies
    - The when users hover over the `Book Stores` menu, all of the above mentioned product categories should be visible. When the user clicks on a category redirect the user to product list page with the product category pre-filtered
    - The `Printing Services` menu will have two sub sections called `Products Offered` and `Inhouse Services`
    - `Products Offered` will have the following sub options: Monocarton, Lid and Tray Box, Rigid Case, Corrugated Box, Books and Pads, Calendars, Pamplets and Brochures, Invitations and Business Cards, Envelopes and Pouches, Files and Folders, Flyers and Danglers
    - `Inhouse Services` will have the following sub options: Offset Printing, Met. PET Printing, UV - Gloss/Matt Coating, Hot Foil Printing, Screen Printing, Food Grade Printing, Rigid Case, Corrugation, Lamination, Pasting, Die Cutting, Embossing/Engraving, Binding and Publishing
4. The Wall of fame can have some dummy content
5. Contact us page will have the address, a map integration, social media links and a form to capture user queries: It should have name, email, phone number, and message
6. The about us page can have dummy content
7. There should be a terms and conditions page
8. There should be a privacy policy page
9. A cart page is required
10. The footer should have the following links: Terms and conditions, About us, Wall of fame links. It will also have links to social media handles for Facebook, instagram, and WhatsApp
11. We need sign up by email & password and facebook login. During sign up we need to verify the email by OTP, and then ask for the password
12. The accounts page will have 3 tabs:
    - My profile: will have first name, last name, phone number, email, and password. The user needs to verify his OTP again if we changes his email
    - My address: An address field will have the following fields: Address Line 1 , Address Line 2, Landmark, City, State, Pincode, Company Name, GST Number, Phone
    - My orders: Should show the order history, it is a paginated list. An order can have 3 status: Open, fulfilled, and cancelled. The User can cancel `open` orders
13. Product details page:
    - Will show the details regarding a single product
    - there will one button to add the product to cart.
    - It should show the product images similar to Amazon with Zoom feature
14. The checkout page should do the necessary tax calculations on the items that are in the cart
    - If the cart is empty the user should not be allowed to access the checkout page
    - The checkout page will have an option to add a coupon code
    - The coupon code can either be flat offer of certain number or percentage of the total bill
    - The coupon can be added using a script to the DB in the server side   
15. The website should be mobile responsive
16. The frontend is a static site generated using NextJS V16. Follow the App Router architecture with route groups, and  Static Site Generation (SSG) with generateStaticParams
17. Need to follow SOLID principles, Abstract Factory, Adapter, Builder, Bridge,  and other relevant design patterns in the server side.
18. Use Zustand V5 plugin for state management in NextJS. Follow Zustand store slice pattern
19. Need a login page for the site administrator so that he can access the admin site
20. The admin site should have a page for CRUD operations for products
21. The admin site should have a page for CRUD operations for categories
22. The admin site should have a page to manage the sales orders
23. Need seed commands / scrips in the server side for the following:
    - Adding two default coupons to the DB: `SUPERHIT`: Flat 150 off, and `FIRSTTIME` 10 percent off.
    - Adding the default admin user into the system
    - Adding 5 dummy products to the system with 3 categories

---

## Required Output

Generate a phased plan with the following structure:

### Phase 0 — Project Setup
Tasks to scaffold the project, configure tooling, set up CI/CD, and establish the development environment.

### Phase 1 — Core Infrastructure
Tasks to build the foundational layers: database schema and migrations, auth system, base API structure, frontend scaffolding and routing.

### Phase 2 — Feature Implementation
One sub-phase per v1 feature. Each sub-phase contains:
- Backend tasks (models, handlers, tests)
- Frontend tasks (components, pages, API client integration)
- Integration tasks

### Phase 3 — Quality & Hardening
Tasks for: error handling, input validation, rate limiting, logging, test coverage gaps, security review.

### Phase 4 — Deployment & Launch
Tasks for: production infrastructure, CI/CD pipeline, monitoring setup, documentation, launch checklist.

---

## Task Format

For each task, use this format:

```
#### TASK-NNN: <Task name>
- **Phase:** <phase name>
- **Depends on:** <TASK-NNN list or "none">
- **Files:** <list of files to create or modify>
- **Description:** <2-3 sentence description of what to do>
- **Acceptance criteria:**
  - [ ] ...
```

---

## Additional Requirements

- Every task must have clear acceptance criteria
- Flag tasks that require an ADR with `> **ADR needed:** <topic>`
- Flag tasks that must be done sequentially vs. in parallel
- Total task count should be between 20 and 50 for a typical v1
- Include a dependency graph summary at the end (which phases can overlap)
