import { TerminalError, workflow, WorkflowContext } from "@restatedev/restate-sdk"
const OrderManager = {
    validate: (order) => {
        console.log("Order validation logic", order)
    },
    placeOnExchange: (order) => {
        console.log("Order placement logic", order)
    },
    reserve: (order) => {
        console.log("Order reservation logic", order)

    },
    settle: (order) => {
        console.log("Order settlement logic", order)

    },
    release: (order) => {
        console.log("Order placement logic", order)

    },

}
const tradeWorkflowName = "tradeWorkflow"
export const tradeWorkflow = workflow({
    name: tradeWorkflowName,
    handlers: {
        run: async (ctx: WorkflowContext, order: any) => {
            const compensations = [];
            //Workflow:

            // Validate order

            // Reserve funds

            // Place an order on the exchange

            // Settle trade
            try {
                const orderId = ctx.rand.uuidv4();
                compensations.push(() =>
                    ctx.run("release-order", () => OrderManager.release(orderId)),
                );
                const orderRef = await ctx.run("validate", () =>
                    OrderManager.validate(order),
                );
                compensations.push(() =>
                    ctx.run(`release-order-${orderId}`, () => OrderManager.release(orderId)),
                );
                const reservedOrder = await ctx.run("reserve-fund", () =>
                    OrderManager.reserve(orderRef),
                );
                compensations.push(() =>
                    ctx.run(`release-order-${orderId}`, () => OrderManager.release(orderId)),
                );
                const placedOrder = await ctx.run("place-order", () =>
                    OrderManager.placeOnExchange(reservedOrder),
                );
                compensations.push(() =>
                    ctx.run(`release-order-${orderId}`, () => OrderManager.release(orderId)),
                );
                const settledTrade = await ctx.run("settle-trade", () =>
                    OrderManager.settle(placedOrder),
                );
                return {
                    data: settledTrade,
                    message: "Order placed successfully"
                }

            } catch (error) {
                if (error instanceof TerminalError) {
                    for (const compensation of compensations.reverse()) {
                        await compensation();
                    }
                }

            }

        },
    },
});

