package routers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Cinnamoon-dev/blue-gopher/internal/domain"
	"github.com/Cinnamoon-dev/blue-gopher/internal/http/handlers"
	"github.com/Cinnamoon-dev/blue-gopher/internal/middleware"
)

type OrderRouter struct {
	DB *sql.DB
}

func NewOrderRouter(db *sql.DB) OrderRouter {
	return OrderRouter{
		DB: db,
	}
}

type OrderRequest struct {
	UserId int64 `json:"user_id"`
	Value  int64 `json:"value"`
}

func (ro *OrderRouter) BaseRoutes() http.HandlerFunc {
	router := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			// Pega o payload - OK!
			// Verifica se o usuário existe - OK!
			// Verifica se existe alguma linha em orders com esse userId, senão crie com valores vazios/0 - OK!
			// Faz um update atualizando o valor das linhas - OK!
			var request OrderRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				handlers.RespondError(w, err)
				return
			}

			row := ro.DB.QueryRow(`
				SELECT id
				FROM usuarios	
				WHERE id = ?
				LIMIT 1
			`, request.UserId)

			var id int64
			if err := row.Scan(&id); err != nil {
				handlers.RespondError(w, err)
				return
			}

			row = ro.DB.QueryRow(`
				SELECT id, user_id, total_spent, order_quantity
				FROM compras
				WHERE user_id = ?	
			`, request.UserId)

			var order domain.Order
			if err := row.Scan(&order.ID, &order.UserID, &order.TotalSpent, &order.OrderQuantity); err != nil {
				fmt.Printf("row.Scan: %s\n", err.Error())
				_, err := ro.DB.Exec(`
					INSERT INTO compras(user_id, total_spent, order_quantity) VALUES (?, ?, ?);
				`, request.UserId, request.Value, 1)

				if err != nil {
					fmt.Printf("Insert into: %s\n", err.Error())
					handlers.RespondError(w, err)
					return
				}

				handlers.RespondJSON(w, http.StatusOK, map[string]string{"message": "Order added successfully"})
				return
			}

			// Como deve existir apenas uma linha nessa tabela por usuário faz sentido filtrar por user_id
			_, err := ro.DB.Exec(`
				UPDATE compras
				SET total_spent = ?, order_quantity = ?
				WHERE user_id = ?;
			`, order.TotalSpent+request.Value, order.OrderQuantity+1, request.UserId)

			if err != nil {
				fmt.Printf("update: %s\n", err.Error())
				handlers.RespondError(w, err)
				return
			}

			handlers.RespondJSON(w, http.StatusOK, map[string]string{"message": "Order added successfully"})
			return
		}
	})

	router = middleware.Logging(router)
	return router
}
