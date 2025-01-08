import { Inject, Injectable } from '@nestjs/common';
import { CreateRouteDto } from './dto/create-route.dto';
import { UpdateRouteDto } from './dto/update-route.dto';
import { PrismaService } from 'src/prisma/prisma.service';
import { DirectionsService } from 'src/maps/directions/directions.service';
import * as kafkaLib from '@confluentinc/kafka-javascript';

@Injectable()
export class RoutesService {
  constructor(
    private prismaService: PrismaService,
    private directionsService: DirectionsService,
    @Inject('KAFKA_PRODUCER') private KafkaProducer: kafkaLib.KafkaJS.Producer,
  ) {}

  async create(CreateRouteDto: CreateRouteDto) {
    const { available_travel_modes, geocoded_waypoints, routes, request } =
      await this.directionsService.getDirections(
        CreateRouteDto.source_id,
        CreateRouteDto.destination_id,
      );

    const legs = routes[0].legs[0];

    const route = await this.prismaService.route.create({
      data: {
        name: CreateRouteDto.name,
        source: {
          name: legs.start_address,
          location: {
            lat: legs.start_location.lat,
            lng: legs.start_location.lng,
          },
        },
        destination: {
          name: legs.end_address,
          location: {
            lat: legs.end_location.lat,
            lng: legs.end_location.lng,
          },
        },
        duration: legs.duration.value,
        distance: legs.distance.value,
        directions: JSON.parse(
          JSON.stringify({
            available_travel_modes,
            geocoded_waypoints,
            routes,
            request,
          }),
        ),
      },
    });

    await this.KafkaProducer.send({
      topic: 'route',
      messages: [
        {
          key: route.id.toString(),
          value: JSON.stringify({
            event: 'route_created',
            id: route.id,
            distance: legs.distance.value,
            directions: legs.steps.reduce((acc, step) => {
              acc.push({
                let: step.start_location.lat,
                lng: step.start_location.lng,
              });
              acc.push({
                let: step.end_location.lat,
                lng: step.end_location.lng,
              });
              return acc;
            }, []),
          }),
        },
      ],
    });

    return route;
  }

  findAll() {
    return this.prismaService.route.findMany();
  }

  findOne(id: string) {
    return this.prismaService.route.findUniqueOrThrow({
      where: { id },
    });
  }

  update(id: string, updateRouteDto: UpdateRouteDto) {
    return this.prismaService.route.update({
      where: { id },
      data: updateRouteDto,
    });
  }

  remove(id: number) {
    return `This action removes a #${id} route`;
  }
}
