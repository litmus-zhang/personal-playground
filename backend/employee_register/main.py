class EmployeeRegister:
    def __init__(self):
        self.workers = {}

    def add_worker(self, worker_id: str, compensation: int, position: str):
        self.workers[worker_id] = {
            "compensation": compensation,
            "position": position,
            "timeIn": [],
            "timeOut": [],
        }

    def register(self, worker_id: str, timestamp: int):
        if worker_id not in self.workers:
            return "Worker not found"
        timeIn = self.workers[worker_id]["timeIn"]
        timeOut = self.workers[worker_id]["timeOut"]
        if len(timeIn) == len(timeOut):
            timeIn.append(timestamp)
        else:
            timeOut.append(timestamp)

        return "Registered"

    def get(self, worker_id: str):
        if worker_id not in self.workers:
            return "Worker not found"
        stats = self.workers[worker_id]
        n = min(len(stats["timeIn"]), len(stats["timeOut"]))
        total_time = 0
        for i in range(n):
            total_time += stats["timeOut"][i] - stats["timeIn"][i]

        return {
            "workerId": worker_id,
            "compensation": stats["compensation"],
            "position": stats["position"],
            "totalTime": total_time,
        }


