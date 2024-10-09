from main import EmployeeRegister


def test_init():
    register = EmployeeRegister()
    assert register.workers == {}


def test_add_worker():
    register = EmployeeRegister()
    register.add_worker("1", 1000, "developer")
    assert register.workers == {
        "1": {
            "compensation": 1000,
            "position": "developer",
            "timeIn": [],
            "timeOut": [],
        }
    }


def test_add_multiple_workers():
    register = EmployeeRegister()
    # use table driven tests to test multiple workers
    workers = [
        ("1", 1000, "developer"),
        ("2", 2000, "manager"),
        ("3", 3000, "designer"),
    ]
    for worker in workers:
        register.add_worker(*worker)
    assert register.workers.keys() == {"1", "2", "3"}


def test_register():
    register = EmployeeRegister()
    register.add_worker("1", 1000, "developer")
    # use table driven tests to test multiple timestamps
    timestamps = [1000, 2000, 3000]
    for timestamp in timestamps:
        assert register.register("1", timestamp) == "Registered"
    assert register.workers["1"]["timeIn"] == [1000, 3000]
    assert register.workers["1"]["timeOut"] == [2000]


def test_get():
    register = EmployeeRegister()
    register.add_worker("1", 1000, "developer")
    timestamps = [1000, 2000, 3000, 4000, 7200, 10230]
    for timestamp in timestamps:
        register.register("1", timestamp)
    assert register.get("1")["totalTime"]== 5030
