# Prompt: Generate Product Requirements Document (PRD)

> **How to use:** Fill in every `[[PLACEHOLDER]]` from your intake questionnaire, then paste this entire file into your AI tool. Save the output to `docs/02_outputs/01_prd.md`.

---

## Instructions to AI

You are a senior product manager. Generate a comprehensive, production-ready **Product Requirements Document (PRD)** based on the project details below.

The PRD must be structured for an engineering team to begin implementation immediately. Be specific and opinionated — do not leave ambiguous sections. Where the input is vague, make reasonable assumptions and note them explicitly.

---

## Project Input

**Project name:** e-commerce-site
**Description:** This is a simple e-commerce-site for selling books
**Team:** Feather Tech

**Primary users:** A simple e-commerce site for selling books online
**Problem being solved:** The product allows customer's to buy books online. It is not a marketplace for buyers and sellers it is a simple e-commerce site where customers can order books and the website owner will dispatch the orders
**Success metrics (6 months):** Don't focus on the success of the project

**Must-have features (v1):**
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
15. Use Zustand V5 plugin for state management in NextJS
16. Need to follow SOLID principles, Abstract Factory, Adapter, Builder, Bridge,  and other relevant design patterns in the server side. Do something similar to the frontend.

**Out of scope (v1):**
1. `Printing Services` page
2. Navbar links for `Printing Services` drop down

**External integrations:**
No

**Target launch date:** no fixed date

---

## Required PRD Sections

Generate the PRD with all of the following sections:

1. **Executive Summary** — 2-3 sentences: what, who, why
2. **Problem Statement** — current pain, root cause, impact
3. **Goals & Success Metrics** — SMART goals with measurable KPIs
4. **Non-Goals** — explicit scope exclusions
5. **User Personas** — 2-3 personas with needs and pain points
6. **User Stories** — written as "As a [persona], I want [action] so that [benefit]" — cover all v1 features
7. **Functional Requirements** — numbered, grouped by feature area
8. **Non-Functional Requirements** — performance, security, scalability, availability
9. **Assumptions** — list every assumption made
10. **Open Questions** — what must be resolved before or during implementation
11. **Out of Scope** — confirm what is explicitly excluded
12. **Timeline** — high-level milestone breakdown toward launch date

---

## Output Format

- Use Markdown with clear `##` and `###` headings
- Number all functional requirements (e.g. FR-001, FR-002)
- Flag assumptions with `> **Assumption:**`
- Flag open questions with `> **Open question:**`
- Target length: 1500–3000 words
