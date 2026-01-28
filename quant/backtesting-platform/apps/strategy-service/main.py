class StrategyService(strategy_pb2_grpc.StrategyServiceServicer):
    def RegisterStrategy(self, request, context):
        strategy_id = str(uuid.uuid4())
        save_strategy(strategy_id, request.source_code)
        return strategy_pb2.RegisterStrategyResponse(
            strategy_id=strategy_id
        )