# type: ignore

from diagrams import Cluster, Diagram, Edge, Node
from diagrams.aws.database import RDS
from diagrams.aws.network import ELB
from diagrams.onprem.inmemory import Redis
from diagrams.programming.language import Go

with Diagram(
    "",
    show=False,
    filename="./assets/architecture",
    graph_attr={"splines": "spline"},
):
    lb = ELB("LoadBalancer")

    with Cluster("App Servers"):
        srv3 = Go("app-3")
        srv2 = Go("app-2")
        srv1 = Go("app-1")

    with Cluster("Queue workers"):
        click3 = Go("click-tracker-3")
        click2 = Go("click-tracker-2")
        click1 = Go("click-tracker-1")

    db = RDS("Postgres")
    redis = Redis("Redis")

    lb >> srv2
    srv2 >> Edge(label="Read/Write URLs") >> db
    srv2 >> Edge(label="Queue click for storage") >> redis

    redis >> click2
    click2 >> db
