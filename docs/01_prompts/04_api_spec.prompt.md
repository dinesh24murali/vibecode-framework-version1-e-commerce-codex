# Prompt: Generate OpenAPI 3.1 Specification

> **How to use:** Complete Prompts 01–03 first. Fill in every `[[PLACEHOLDER]]`, then paste this entire file into your AI tool. Save the output to `docs/02_outputs/04_api_spec.yaml`. This file becomes the contract for `tests/contract/`.

---

## Instructions to AI

You are a senior API designer. Generate a complete **OpenAPI 3.1 specification** in YAML format for the project below.

This spec is the **single source of truth** for the API. The frontend client will be generated from it. Contract tests in `tests/contract/` will validate against it. Do not leave placeholder paths — generate real, complete endpoint definitions.

---

## Project Input

**Project name:** e-commerce-site
**Backend:** Go 1.26 / Gin Web Framework
**Auth strategy:** JWT
**Database:** Postgres

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

**User roles:**
There are two roles in the system
   - customer
   - admin

---

## Requirements

Generate a complete OpenAPI 3.1 YAML spec with:

1. **Info block** — title, version `1.0.0`, description
2. **Servers** — `http://localhost:8080` for dev; add staging/prod stubs
3. **Security schemes** — matching `JWT` (e.g. BearerAuth for JWT)
4. **Tags** — one per feature area
5. **Paths** — for every v1 feature, define:
   - All CRUD endpoints that make sense
   - Request body schema (inline `$ref` components)
   - Response schemas for 200, 201, 400, 401, 403, 404, 422, 500
   - `security` field on protected endpoints
6. **Components / Schemas** — reusable schemas for all entities
7. **Components / Parameters** — reusable path/query params (e.g. `id`, `page`, `limit`)
8. **Components / Responses** — reusable error response shapes
9. **Pagination** — use cursor-based or offset pagination consistently; define once, reference everywhere

---

## Output Format

- Output **only valid YAML** — no markdown code fences, no explanation text
- Start with `openapi: "3.1.0"`
- Use `$ref` for all reusable schemas — do not inline complex objects
- All schema properties must have `type` and `description`
- Include at least one example per schema
