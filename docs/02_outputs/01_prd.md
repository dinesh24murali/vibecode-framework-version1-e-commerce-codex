## Executive Summary

`e-commerce-site` is a direct-to-consumer online bookstore for Feather Tech that allows customers to discover, filter, and order books and related paper products online while the business owner fulfills orders manually. The v1 release focuses on a fast mobile-responsive storefront, account and order management, coupon-enabled checkout without online payments, and a lightweight admin portal for catalog and order operations.

## Problem Statement

The business currently lacks a structured online commerce experience for selling books and related stationery products directly to customers. Customers cannot easily browse categories, compare products, search by attributes such as level or theme, or place orders with a persistent account and order history. This gap exists because there is no integrated storefront, customer account system, or internal admin tooling for managing products, categories, and sales orders. The impact is lost sales, inefficient manual order capture, inconsistent customer communication, and no reliable operational record of catalog and order data.

## Goals & Success Metrics

The product goal is to launch a production-ready v1 bookstore that supports catalog discovery, cart, checkout, customer accounts, and admin operations with enough structure for immediate engineering execution.

1. Launch a complete v1 storefront and admin experience that covers 100% of the must-have features listed in this PRD, excluding explicitly documented non-goals, before public release.
2. Ensure customers can complete the order flow from landing page to order placement with a median happy-path completion time of under 5 minutes on mobile and desktop.
3. Achieve at least 95% successful completion for email OTP verification, login, add-to-cart, and order placement flows during UAT.
4. Keep product list and product detail pages performant, with Largest Contentful Paint under 2.5 seconds on a 4G mobile connection for cached static pages and under 3.5 seconds for first uncached load.
5. Maintain operational accuracy with 0 critical defects in catalog CRUD, coupon application, tax calculation, and order status transitions at launch.

> **Assumption:** The business does not want the PRD to define commercial growth targets for the first 6 months, so the KPIs focus on product completeness, usability, and operational reliability rather than revenue or conversion targets.

## Non-Goals

- Online payment gateway integration and payment reconciliation
- Marketplace functionality with third-party sellers
- A production `Printing Services` destination page in v1
- Complex promotion engines beyond manual coupon code support
- Customer reviews, ratings, wishlists, and recommendations
- Advanced admin analytics, inventory forecasting, or shipment tracking integrations
- Multi-language, multi-currency, or international tax support in v1

## User Personas

### Persona 1: Retail Book Buyer

An end customer browsing on mobile or desktop to purchase books and stationery for personal use or gifting.

Needs:
- Quick product discovery through navigation, search, and filters
- Clear product images, pricing, discount, stock status, and order total
- A simple account flow with saved addresses and order history

Pain points:
- Difficulty finding relevant products across many categories
- Friction in checkout and address entry
- Lack of clarity on order status after purchase

### Persona 2: Returning Account Holder

A repeat customer who wants to reorder, manage profile details, maintain addresses, and review prior orders.

Needs:
- Secure login with reliable account recovery and email verification
- Editable profile and addresses
- Ability to view and cancel open orders

Pain points:
- Re-entering details repeatedly
- Unclear history of previous purchases
- Friction when updating email or contact information

### Persona 3: Store Administrator

An internal operator responsible for maintaining products, categories, coupons, and sales orders.

Needs:
- Secure admin login
- Simple CRUD workflows for catalog management
- Clear order statuses and operational order management

Pain points:
- Manual spreadsheet-based catalog updates
- No centralized order queue
- Inconsistent product metadata and stock visibility

## User Stories

- As a retail book buyer, I want to land on a visually engaging homepage with a carousel and product cards so that I can quickly discover featured products.
- As a retail book buyer, I want to browse products on a listing page with left-side filters so that I can narrow the catalog by category, level, theme, and size.
- As a retail book buyer, I want search and filter state preserved in the URL so that I can share or revisit the same results.
- As a retail book buyer, I want to search for products from the navbar so that I can reach relevant products from any page.
- As a retail book buyer, I want the Book Store menu to expose all product categories so that I can jump directly to a pre-filtered catalog page.
- As a retail book buyer, I want to view a product detail page with zoomable images so that I can inspect the product before purchase.
- As a retail book buyer, I want to add products to my cart from the listing page and product detail page so that I can build an order quickly.
- As a retail book buyer, I want to review my cart contents before checkout so that I can confirm quantities and pricing.
- As a retail book buyer, I want checkout to calculate taxes and coupon discounts so that I know the exact payable amount before placing the order.
- As a retail book buyer, I want to be blocked from checkout when my cart is empty so that I do not enter a broken flow.
- As a new customer, I want to sign up with email OTP verification and Facebook login so that I can create an account securely and conveniently.
- As a returning account holder, I want to manage my profile, addresses, and password so that my account data stays current.
- As a returning account holder, I want to re-verify by OTP when changing my email so that my account stays secure.
- As a returning account holder, I want to view paginated order history and cancel open orders so that I can manage active purchases.
- As a visitor, I want to view About Us, Contact Us, Terms and Conditions, Privacy Policy, and Wall of Fame pages so that I can trust the business and contact it if needed.
- As a visitor, I want the Contact Us page to include address, map, social links, and inquiry form so that I can reach the business through my preferred channel.
- As a store administrator, I want a secure admin login so that only authorized staff can access back-office functions.
- As a store administrator, I want CRUD pages for products and categories so that I can keep the catalog current.
- As a store administrator, I want to manage sales orders so that I can process, fulfill, or cancel them accurately.
- As an operator, I want seed scripts for coupons, admin user, and sample products so that local and staging environments can be initialized quickly.

## Functional Requirements

### Storefront and Navigation

- FR-001: The system shall provide a public landing page with a horizontal hero carousel followed by a featured product card grid.
- FR-002: The global navbar shall appear on all public pages and include links for `Book Store`, `Wall of fame`, `Contact us`, `Printing service`, and `Printing supplies`.
- FR-003: Hovering or tapping `Book Store` shall reveal all supported product categories: Colouring Books, Drawing Books, Journals, Calendars, Sketch Pads, Planners, Stained Art, Fabrii Art, Note Book, Sticker Book, Sticker Colouring Book, and Activity Cards.
- FR-004: Selecting a category from the `Book Store` menu shall route the user to the product listing page with the chosen category applied as a URL query parameter.
- FR-005: The navbar shall include a search input on the top right and submit users to the product listing page with the search term preserved in the URL.
- FR-006: The footer shall appear on all public pages and include links for `Terms and conditions`, `About us`, and `Wall of fame`, plus Facebook, Instagram, and WhatsApp social links.

### Product Catalog

- FR-007: The product listing page shall display products in a paginated grid with product image, name, price, discount, stock status, and add-to-cart action.
- FR-008: Every product shall store and expose these attributes: category, level, theme, size, price, discount, stock quantity or stock status, description, and images.
- FR-009: The product listing page shall render a left-side filter panel with category, level, theme, and size filters.
- FR-010: The product listing page shall allow combining multiple filters with a search term.
- FR-011: The current product listing state, including search, category, level, theme, size, sort, and page number, shall be maintained in the URL query string.
- FR-012: The product detail page shall display product metadata, price, discount, stock status, description, and a gallery with thumbnail selection and image zoom behavior similar to mainstream e-commerce sites.
- FR-013: Users shall be able to add a product to cart from both listing and detail pages.

### Cart and Checkout

- FR-014: The cart page shall display selected products, quantities, unit price, discount, subtotal, tax, coupon discount, and final total.
- FR-015: Users shall be able to update item quantities and remove items from the cart.
- FR-016: The checkout page shall be inaccessible when the cart is empty and shall redirect the user back to the cart with a visible message.
- FR-017: Checkout shall collect or confirm a shipping address from the authenticated user's saved addresses.
- FR-018: Checkout shall support applying one coupon code per order.
- FR-019: The coupon engine shall support both flat amount discounts and percentage-based discounts.
- FR-020: Coupon definitions shall be created and updated server-side, including via seed scripts.
- FR-021: Checkout shall calculate tax, coupon discount, and final total deterministically on the server before order placement.
- FR-022: Because online payment is out of scope, checkout shall place orders in an unpaid manual-fulfillment state using a cash-on-delivery or pay-later style order record.

### Authentication and Accounts

- FR-023: The system shall support customer signup with email and password.
- FR-024: The email signup flow shall verify email ownership via OTP before the customer sets a password and activates the account.
- FR-025: The system shall support Facebook social login for customers.
- FR-026: The customer account area shall include `My profile`, `My address`, and `My orders` tabs.
- FR-027: `My profile` shall support first name, last name, phone number, email, and password management.
- FR-028: When a customer changes their email address, the system shall require OTP verification before the new email becomes active.
- FR-029: `My address` shall allow create, read, update, and delete for addresses with Address Line 1, Address Line 2, Landmark, City, State, Pincode, Company Name, GST Number, and Phone.
- FR-030: `My orders` shall display paginated order history with status values limited to `Open`, `Fulfilled`, and `Cancelled`.
- FR-031: Customers shall be able to cancel only orders in `Open` status.

### Content Pages and Contact

- FR-032: The system shall provide public pages for About Us, Terms and Conditions, Privacy Policy, Wall of Fame, and Contact Us.
- FR-033: The Wall of Fame and About Us pages may launch with approved placeholder content.
- FR-034: The Contact Us page shall display business address, an embedded map, Facebook/Instagram/WhatsApp links, and an inquiry form with name, email, phone number, and message.
- FR-035: Contact form submissions shall be stored server-side for admin review or emailed to a configured business inbox.

### Admin Portal

- FR-036: The system shall provide a separate admin login experience for authorized administrators.
- FR-037: The admin portal shall provide product CRUD including images, pricing, discount, stock, and the custom attributes category, level, theme, and size.
- FR-038: The admin portal shall provide category CRUD for the supported product taxonomy.
- FR-039: The admin portal shall provide a sales order management page with order list, customer details, order totals, and status management.
- FR-040: Admins shall be able to change order status between `Open`, `Fulfilled`, and `Cancelled`, with auditability of who made the change.

### Seeds, Architecture, and Platform Constraints

- FR-041: The backend shall include seed commands or scripts that create coupon `SUPERHIT` with flat 150 discount, coupon `FIRSTTIME` with 10% discount, a default admin user, and 5 dummy products spanning 3 categories.
- FR-042: The frontend shall be implemented with NextJS 16 App Router using route groups and static site generation where feasible, including `generateStaticParams` for statically generated routes.
- FR-043: The frontend shall use Zustand v5 with a slice-based store structure for cart, auth, and UI state.
- FR-044: The backend shall follow SOLID principles and use relevant patterns such as Abstract Factory, Adapter, Builder, and Bridge where they provide clear maintainability value.
- FR-045: The site shall be responsive across mobile, tablet, and desktop breakpoints.

## Non-Functional Requirements

- Performance: Public marketing and catalog pages shall be optimized for SSG-first delivery, compressed images, lazy loading below the fold, and sub-2.5 second LCP targets on key pages.
- Security: Passwords shall be hashed with a modern algorithm, OTPs shall expire within 10 minutes, sessions or JWTs shall be protected against replay and CSRF risks, and admin routes shall enforce role-based authorization.
- Privacy: Contact form and account data shall be stored securely, transmitted over HTTPS only, and handled in line with the published privacy policy.
- Availability: The system shall target 99.5% monthly uptime for public storefront and admin access excluding scheduled maintenance.
- Scalability: The architecture shall support at least 10,000 product records and 200 concurrent users without redesign.
- Reliability: Tax, discount, and order total calculations shall be server-authoritative and idempotent for repeated checkout submissions.
- Maintainability: API contracts shall remain aligned with the OpenAPI spec, and frontend API access shall use generated clients rather than hand-written duplicates.
- Accessibility: Public pages and core forms shall meet WCAG 2.1 AA expectations for keyboard navigation, color contrast, labels, and form error messaging.
- Observability: The platform shall log authentication events, admin CRUD actions, contact form submissions, and order status changes with timestamp and actor context.

## Assumptions

> **Assumption:** The site serves a single business owner and a single catalog, not multiple stores or sellers.

> **Assumption:** The primary operating geography is India because address requirements include GST Number and Pincode; therefore tax logic should be designed around GST-style calculations unless clarified otherwise.

> **Assumption:** v1 supports manual order placement without online payment, and orders enter the system as pending fulfillment rather than paid.

> **Assumption:** Guest browsing is allowed, but placing an order requires user authentication so orders can be tied to an account and address book.

> **Assumption:** `Library` in the footer links to the main product listing or storefront collection page.

> **Assumption:** `Printing service` and `Printing supplies` remain visible as navigation labels for information architecture consistency, but their dropdown workflows and destination page are out of scope for v1.

> **Assumption:** Search is keyword-based across product name, category, theme, and description, not semantic or typo-tolerant in v1.

> **Assumption:** Product inventory in v1 is managed as a simple stock count or in-stock/out-of-stock state without warehouse-level tracking.

> **Assumption:** The admin portal is part of the same overall product stack and deployment, protected by admin-specific authentication and authorization.

> **Assumption:** OTP delivery will be implemented through email only; SMS OTP is not included in v1.

## Open Questions

> **Open question:** What exact tax rules should be applied at checkout: a single GST percentage for all products, category-specific rates, or state-aware tax treatment?
- The tax is at individual product level. Need to calculate the tax based in the rate mentioned for the individual product

> **Open question:** Should customers be allowed to maintain multiple saved addresses and choose one default address at checkout?
- Yes customers can maintain multiple saved addresses and choose one default address

> **Open question:** What is the expected fulfillment communication after order placement if no payment gateway exists: confirmation email only, manual phone follow-up, or both?
- Both confirmation email and Manual phone follow-up

> **Open question:** What content and destination should `Printing supplies` represent in v1, given it is listed in the navbar but not otherwise specified?
- You can skip the content for now. Leave a blank page page

> **Open question:** Should the contact form create admin-visible records, send transactional email, or both?
- Both

> **Open question:** What fields are mandatory versus optional for product creation in admin, especially for level, theme, and size?
- all fields are mandatory

> **Open question:** Is Facebook login mandatory at launch if app credentials or platform approval are delayed, or may it be feature-flagged behind email signup?
- A feature-flagged can be set that will disable Facebook login until platform approval

> **Open question:** Should order cancellation by the customer require a cancellation reason?
- yes and optional cancellation reason

## Out of Scope

- `Printing Services` page implementation
- Navbar dropdown implementation details for `Printing Services`
- Payment gateway integration, card processing, or wallet support
- Shipment tracking, courier integration, and automated dispatch workflows
- Reviews, ratings, wishlists, loyalty points, and referrals
- Marketplace seller onboarding or multi-vendor workflows

## Timeline

### Milestone 1: Product Definition and Architecture

Week 1:
- Finalize PRD, assumptions, and unresolved product questions
- Produce architecture document, OpenAPI draft, and initial ADRs
- Confirm auth, tax, coupon, and admin access model

### Milestone 2: Backend Foundations

Weeks 2-4:
- Set up data models, migrations, auth flows, and admin authorization
- Implement catalog, cart, checkout, coupon, address, and order APIs
- Build seed scripts for coupons, admin user, and dummy products

### Milestone 3: Frontend Storefront

Weeks 3-6:
- Implement App Router structure, shared layout, navbar, footer, and static public pages
- Build landing, product listing, product detail, cart, checkout, and account flows
- Add Zustand store slices and API client integration

### Milestone 4: Admin Portal and Hardening

Weeks 6-7:
- Implement admin login, product CRUD, category CRUD, and sales order management
- Complete responsive behavior, accessibility fixes, and QA for core journeys
- Execute verification scripts, regression tests, and launch readiness review

### Milestone 5: Launch Preparation

Week 8:
- Seed production-like environment
- Validate SEO, analytics hooks if required, and operational runbooks
- Resolve UAT defects and approve launch
