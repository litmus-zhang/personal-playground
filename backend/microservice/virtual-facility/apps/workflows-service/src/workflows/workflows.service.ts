import { Injectable, Logger, NotFoundException } from '@nestjs/common';
import { CreateWorkflowDto, UpdateWorkflowDto } from '@app/workflows';
import { Workflow } from './entities/workflow.entity';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';

@Injectable()
export class WorkflowsService {
  private readonly logger = new Logger(WorkflowsService.name);
  constructor(
    @InjectRepository(Workflow)
    private readonly workflowRepository: Repository<Workflow>,
  ) {}

  async create(createWorkflowDto: CreateWorkflowDto): Promise<Workflow> {
    const workflow = this.workflowRepository.create({
      name: createWorkflowDto.name,
      buildingId: createWorkflowDto.buildingId,
    });

    const newWorkflowEntity = await this.workflowRepository.save(workflow);
    this.logger.debug(
      `Created workflow with id ${newWorkflowEntity.id} for building ${newWorkflowEntity.buildingId}`,
    );
    return newWorkflowEntity;
  }

  findAll(): Promise<Workflow[]> {
    return this.workflowRepository.find();
  }

  async findOne(id: number): Promise<Workflow> {
    const workflow = await this.workflowRepository.findOne({ where: { id } });
    if (!workflow) {
      throw new NotFoundException(`workflow #${id} does not exist`);
    }
    return workflow;
  }

  async update(
    id: number,
    updateWorkflowDto: UpdateWorkflowDto,
  ): Promise<Workflow> {
    const workflow = await this.workflowRepository.preload({
      id: +id,
      name: updateWorkflowDto.name,
      buildingId: updateWorkflowDto.buildingId,
    });
    if (!workflow) {
      throw new NotFoundException(`Workflow #${id} does not exist`);
    }
    return this.workflowRepository.save(workflow);
  }

  async remove(id: number): Promise<Workflow> {
    const workflow = await this.findOne(id);
    return this.workflowRepository.remove(workflow);
  }
}
