# Intake Questionnaire

> Fill this out before using any prompt in `docs/01_prompts/`.
> Your answers replace the `[[PLACEHOLDER]]` tokens in each prompt.

---

## 1. Project Basics

**1.1 Project name:**
`e-commerce-site`

**1.2 One-sentence description:**
`This is a simple e-commerce-site for selling books`

**1.3 Team / company name:**
`Feather Tech`

**1.4 Target launch date (or "no fixed date"):**
no fixed date

---

## 2. Users & Problem

**2.1 Who are the primary users?**
A simple e-commerce site for selling books online

**2.2 What problem does this solve?**
The product allows customer's to buy books online. It is not a marketplace for buyers and sellers it is a simple e-commerce site where customers can order books and the website owner will dispatch the orders

**2.3 What does success look like in 6 months?**
Don't focus on the success of the project


---

## 3. Core Features

**3.1 List the 3–5 must-have features for v1:**
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
10. The footer should have the following links: Library, Terms and conditions, About us, Wall of fame links. It will also have links to social media handles for Facebook, instagram, and WhatsApp
11. We need sign up by email & password and facebook login. During sign up we need to verify the email by OTP, and then ask for the password
12. The accounts page will have 3 tabs:
    - My profile: will have first name, last name, phone number, email, and password. The user needs to verify his OTP again if we changes his email
    - My address: An address field will have the following fields: Address Line 1 , Address Line 2, Landmark, City, State, Pincode, Company Name, GST Number, Phone
    - My orders: Should show the order history, it is a paginated list. An order can have 3 status: Open, fulfilled, and cancelled. The User can cancel `open` orders
13. Product details page:
    - Will show the details regarding a single product
    - there will one button to add the product to cart.
    - It should show the product images similar to Amazon with Zoom feature
14. The website should be mobile responsive
15. Use Zustand plugin for state management in NextJS
16. Need to follow SOLID principles and other relevant design patterns in the server side. Do something similar to the frontend.

**3.2 List features explicitly OUT of scope for v1:**
1. `Printing Services` page
2. Navbar links for `Printing Services` drop down

**3.3 Are there any existing systems this must integrate with?**
No

---

## 4. Tech Stack

**4.1 Backend language:**
`Go 1.26`  (e.g. Go, Python, Node.js, Ruby)

**4.2 Backend framework:**
`Gin Web Framework`  (e.g. Gin, FastAPI, Express, Rails)

**4.3 Frontend framework:**
`NextJS 16`  (e.g. React, Vue, SvelteKit, Next.js)

**4.4 Database:**
`Postgres`  (e.g. PostgreSQL, MySQL, SQLite, MongoDB)

**4.5 Cache / queue (if any):**
Redis

**4.6 Any preferred UI component library?**
shadcn/ui, Tailwind

---

## 5. Authentication & Authorization

**5.1 Auth strategy:**
`JWT`  (e.g. JWT, OAuth2 + JWT, Session cookies, Passkeys)

**5.2 OAuth providers (if any):**
none

**5.3 Authorization model:**
Simple owner-based. The User cannot access the cart, and my accounts pages without signing in.

---

## 6. Deployment & Infrastructure

**6.1 Deployment target:**
`AWS`

**6.2 Containerized:**
Yes — Docker Compose

**6.3 CI/CD preference:**
none

**6.4 Environments needed:**
dev

---

## 7. Scale & Non-Functional Requirements

**7.1 Expected users at launch:**
1000

**7.2 Expected peak requests/second:**
10

**7.3 Data sensitivity:**
(e.g. PII, financial data, public data only)

**7.4 Compliance requirements:**
DPDPA - Digital Personal Data Protection Act, India

**7.5 Uptime target:**
best-effort

---

## 8. Testing Expectations

**8.1 Minimum unit test coverage target:**
70%

**8.2 E2e test tool preference:**
Playwright

**8.3 Contract testing:**
Yes — generate from OpenAPI spec

---

## 9. Open Questions

None
