import { IsString, IsNumber } from 'class-validator';

export class CreateWorkflowDto {
  @IsString()
  name: string;

  @IsNumber()
  buildingId: number;
}
