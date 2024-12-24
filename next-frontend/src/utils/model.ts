/* eslint-disable @typescript-eslint/no-explicit-any */
import { DirectionsResponseData } from "@googlemaps/google-maps-services-js";

export type DirectionsData = DirectionsResponseData & { request: any };
