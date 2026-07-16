async function reserveSeat(showId, seatId) {
  return apiFetch(`/bookings/reserve`, {
    method: "POST",
    body: JSON.stringify({ show_id: showId, seat_id: seatId }),
  });
}

async function getReservedSlots(showId) {
  return apiFetch(`/shows/${showId}/reserved-slots`);
}

async function getCheckoutSummary(showId) {
  return apiFetch(`/bookings/checkout?show_id=${showId}`);
}

async function initiateCheckout(showId) {
  return apiFetch(`/bookings/checkout`, {
    method: "POST",
    body: JSON.stringify({ show_id: showId }),
  });
}