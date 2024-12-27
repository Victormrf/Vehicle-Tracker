"use client";

import { use, useEffect, useRef } from "react";
import { useMap } from "../../hooks/useMap";
import { socket } from "@/utils/socket-io";

export type MapDriverPops = {
  route_id: string | null;
};

export function MapDriver(props: MapDriverPops) {
  const { route_id } = props;
  const mapContainerRef = useRef<HTMLDivElement>(null!);
  const map = useMap(mapContainerRef);

  useEffect(() => {
    if (!map || !route_id) return;

    socket.connect();

    socket.on(`server:new-points/${route_id}:list`, (data) => {});
  }, [route_id, map]);

  return <div className="w-2/3 h-full" ref={mapContainerRef} />;
}
