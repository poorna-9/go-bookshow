// Toggle a seat: selects it if free, deselects it if you already held it.
// Returns { action: "selected" | "unselected", session_id }
async function reserveSeat(showId, seatId) {
  return apiFetch(`/bookings/reserve`, {
    method: "POST",
    body: JSON.stringify({
      user_id: getUserId(),
      show_id: showId,
      seat_id: seatId,
    }),
  });
}

// Full seat map for a show: which seats are available, which are yours,
// which belong to other shoppers, which are already booked.
async function getReservedSlots(showId) {
  return apiFetch(`/shows/${showId}/reserved-slots?user_id=${getUserId()}`);
}

// Preview of your current selection before paying.
async function getCheckoutSummary(showId) {
  return apiFetch(`/bookings/checkout?show_id=${showId}&user_id=${getUserId()}`);
}

// Pay Now: blocks the seats and opens (or reuses) a Razorpay order.
async function initiateCheckout(showId) {
  return apiFetch(`/bookings/checkout`, {
    method: "POST",
    body: JSON.stringify({ user_id: getUserId(), show_id: showId }),
  });
}

// Called after Razorpay's widget finishes, whether it succeeded or was cancelled.
async function paymentCallback(payload) {
  return apiFetch(`/bookings/payment-callback`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
}