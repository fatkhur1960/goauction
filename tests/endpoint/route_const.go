package endpoint

const (
	// AuthorizeUserEndpoint for testing only
	AuthorizeUser = "/auth/v1/authorize"
	// UnauthorizeUserEndpoint for testing only
	UnauthorizeUser = "/auth/v1/unauthorize"
	// CreateChatRoomEndpoint for testing only
	CreateChatRoom = "/chat/v1/new-room"
	// AddProductEndpoint for testing only
	AddProduct = "/product/v1/add"
	// ListProductEndpoint for testing only
	ListProduct = "/product/v1/list"
	// ListMyProductEndpoint for testing only
	ListMyProduct = "/product/v1/me/list"
	// DetailProductEndpoint for testing only
	DetailProduct = "/product/v1/detail/:id"
	// UpdateProductEndpoint for testing only
	UpdateProduct = "/product/v1/update"
	// DeleteProductEndpoint for testing only
	DeleteProduct = "/product/v1/delete/:id"
	// BidProductEndpoint for testing only
	BidProduct = "/product/v1/bid"
	// ReOpenProductBidEndpoint for testing only
	ReOpenProductBid = "/product/v1/reopen"
	// MarkProductAsSoldEndpoint for testing only
	MarkProductAsSold = "/product/v1/mark-as-sold/:id"
	// RegisterUserEndpoint for testing only
	RegisterUser = "/user/v1/register"
	// ActivateUserEndpoint for testing only
	ActivateUser = "/user/v1/activate"
	// MeInfoEndpoint for testing only
	MeInfo = "/user/v1/me/info"
	// UpdateUserInfoEndpoint for testing only
	UpdateUserInfo = "/user/v1/me/info"
	// ListUserBidsEndpoint for testing only
	ListUserBids = "/user/v1/bids"
	// ListUserNotifsEndpoint for testing only
	ListUserNotifs = "/user/v1/notifs"
	// MarkAsReadNotifEndpoint for testing only
	MarkAsReadNotif = "/user/v1/notifs/read"
)