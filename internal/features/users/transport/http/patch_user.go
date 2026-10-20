package users_transport_http

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/darya-kznv/todo-service/internal/core/domain"
	core_errors "github.com/darya-kznv/todo-service/internal/core/errors"
	core_logger "github.com/darya-kznv/todo-service/internal/core/logger"
	core_http_request "github.com/darya-kznv/todo-service/internal/core/transport/http/request"
	core_http_response "github.com/darya-kznv/todo-service/internal/core/transport/http/response"
	core_http_types "github.com/darya-kznv/todo-service/internal/core/transport/http/types"
)

var patchPhoneRegex = regexp.MustCompile(`^\+[0-9]+$`)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name" swaggertype:"string" example:"Maxim Maximovich"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+71112223344"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("'FullName' can't be NULL: %w", core_errors.ErrInvalidArgument)
		}

		fullNameLength := len([]rune(*r.FullName.Value))
		if fullNameLength < 3 || fullNameLength > 100 {
			return fmt.Errorf(
				"`FullName` must be between 3 and 100 symbols: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			if !patchPhoneRegex.MatchString(*r.PhoneNumber.Value) {
				return fmt.Errorf(
					"`PhoneNumber` must start with '+' and contain only digits: %w",
					core_errors.ErrInvalidArgument,
				)
			}

			PhoneNumberLength := len([]rune(*r.PhoneNumber.Value))
			if PhoneNumberLength < 10 || PhoneNumberLength > 15 {
				return fmt.Errorf(
					"`PhoneNumber` must be between 10 and 15 symbols: %w",
					core_errors.ErrInvalidArgument,
				)
			}

			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf(
					"`PhoneNumber` must start with '+' symbol: %w",
					core_errors.ErrInvalidArgument,
				)
			}
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

// PatchUser godoc
// @Summary Обновить пользователя
// @Description Частично обновить пользователя по идентификатору
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID пользователя" Format(uuid)
// @Param request body PatchUserRequest true "Изменяемые поля"
// @Success 200 {object} PatchUserResponse "Обновлённый пользователь"
// @Failure 400 {object} core_http_response.ErrorResponse "Некорректный запрос"
// @Failure 404 {object} core_http_response.ErrorResponse "Пользователь не найден"
// @Failure 409 {object} core_http_response.ErrorResponse "Конфликт версий"
// @Failure 500 {object} core_http_response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /users/{id} [patch]
func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)

		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)

		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FullName.ToDomain(),
		request.PhoneNumber.ToDomain(),
	)
}
