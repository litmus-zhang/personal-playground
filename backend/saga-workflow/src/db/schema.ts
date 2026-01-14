import { sqliteTable, text  } from "drizzle-orm/sqlite-core"
export const order = sqliteTable("order", {
    id: text("id").primaryKey(),
    type: text("type"),
})