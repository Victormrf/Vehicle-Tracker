"use server";

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export async function createRouteAction(state: any, formData: FormData) {
  const { sourceId, destinationId } = Object.fromEntries(formData);

  const directionsResponse = await fetch(
    `http://localhost:3000/directions?originId=${sourceId}&destinationId=${destinationId}`
  );

  if (!directionsResponse.ok) {
    console.error(await directionsResponse.text());
    return { error: "Failed to fetch directions" };
  }

  const directionsData = await directionsResponse.json();

  const { start_address: startAddress, end_address: endAddress } =
    directionsData.routes[0].legs[0];

  const response = await fetch("http://localhost:3000/routes", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      name: `${startAddress} - ${endAddress}`,
      source_id: directionsData.request.origin.place_id.replace(
        "place_id:",
        ""
      ),
      destination_id: directionsData.request.origin.destination_id.replace(
        "place_id:",
        ""
      ),
    }),
  });

  if (!response.ok) {
    console.error(await response.text());
    return { error: "Failed to create route" };
  }

  return { success: true };
}
