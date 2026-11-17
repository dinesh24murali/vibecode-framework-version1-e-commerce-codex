-- name: ListActiveCarouselSlides :many
SELECT * FROM carousel_slides
WHERE is_active = true
ORDER BY sort_order, created_at
LIMIT 5;
