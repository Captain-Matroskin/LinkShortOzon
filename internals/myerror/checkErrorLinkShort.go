package myerror

import (
	"encoding/json"
	"net/http"
)

func (c *CheckError) CheckErrorCreateLinkShort(err error) (error, []byte, int) {
	if err != nil {
		switch err.Error() {
		case LSHCreateLinkShortNotInsertUnique, LSHCreateLinkShortExistsRedis:
			result, errMarshal := json.Marshal(ResultError{
				Status:  http.StatusConflict,
				Explain: err.Error(),
			})
			if errMarshal != nil {
				c.Logger.Errorf("%s, %v, requestId: %d", ErrMarshal, errMarshal, c.RequestId)
				return &MyErrors{
						Text: ErrMarshal,
					},
					nil, http.StatusInternalServerError
			}
			c.Logger.Warnf("%s, requestId: %d", err.Error(), c.RequestId)
			return &MyErrors{
					Text: ErrCheck,
				},
				result, http.StatusOK

		default:
			result, errMarshal := json.Marshal(ResultError{
				Status:  http.StatusInternalServerError,
				Explain: ErrDB,
			})
			if errMarshal != nil {
				c.Logger.Errorf("%s, %v, requestId: %d", ErrMarshal, errMarshal, c.RequestId)
				return &MyErrors{
						Text: ErrMarshal,
					},
					nil, http.StatusInternalServerError
			}
			c.Logger.Errorf("%s, requestId: %d", err.Error(), c.RequestId)
			return &MyErrors{
					Text: ErrCheck,
				},
				result, http.StatusInternalServerError

		}

	}
	return nil, nil, IntNil
}

func (c *CheckError) CheckErrorTakeLinkShort(err error) (error, []byte, int) {
	if err != nil {
		switch err.Error() {
		case LSHTakeLinkShortNotFound:
			result, errMarshal := json.Marshal(ResultError{
				Status:  http.StatusNotFound,
				Explain: LSHTakeLinkShortNotFound,
			})
			if errMarshal != nil {
				c.Logger.Errorf("%s, %v, requestId: %d", ErrMarshal, errMarshal, c.RequestId)
				return &MyErrors{
						Text: ErrMarshal,
					},
					nil, http.StatusInternalServerError
			}
			c.Logger.Warnf("%s, requestId: %d", err.Error(), c.RequestId)
			return &MyErrors{
					Text: ErrCheck,
				},
				result, http.StatusOK

		default:
			result, errMarshal := json.Marshal(ResultError{
				Status:  http.StatusInternalServerError,
				Explain: ErrDB,
			})
			if errMarshal != nil {
				c.Logger.Errorf("%s, %v, requestId: %d", ErrMarshal, errMarshal, c.RequestId)
				return &MyErrors{
						Text: ErrMarshal,
					},
					nil, http.StatusInternalServerError
			}
			c.Logger.Errorf("%s, requestId: %d", err.Error(), c.RequestId)
			return &MyErrors{
					Text: ErrCheck,
				},
				result, http.StatusInternalServerError

		}

	}
	return nil, nil, IntNil
}

func (c *CheckError) CheckErrorCreateLinkShortGrpc(err error) (error, string, int) {
	if err != nil {
		switch err.Error() {
		case LSHCreateLinkShortNotInsertUnique:
			c.Logger.Warnf("%s, requestId: %d", LSHCreateLinkShortNotInsertUnique, c.RequestId)
			return &MyErrors{
					Text: ErrCheck,
				},
				LSHCreateLinkShortNotInsertUnique, http.StatusConflict
		default:
			c.Logger.Errorf("%s, requestId: %d", err.Error(), c.RequestId)
			return &MyErrors{
					Text: ErrInternal,
				},
				ErrDB, http.StatusInternalServerError
		}
	}
	return nil, "", IntNil
}

func (c *CheckError) CheckErrorTakeLinkFullGrpc(err error) (error, string, int) {
	if err != nil {
		switch err.Error() {
		case LSHTakeLinkShortNotFound:
			c.Logger.Warnf("%s, requestId: %d", LSHTakeLinkShortNotFound, c.RequestId)
			return &MyErrors{
					Text: ErrCheck,
				},
				LSHTakeLinkShortNotFound, http.StatusNotFound
		default:
			c.Logger.Errorf("%s, requestId: %d", err.Error(), c.RequestId)
			return &MyErrors{
					Text: ErrInternal,
				},
				ErrDB, http.StatusInternalServerError
		}
	}
	return nil, "", IntNil
}
