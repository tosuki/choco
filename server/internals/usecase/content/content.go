package content

import (
	"choco/server/internals/adapters"
	"choco/server/internals/models"
	"choco/server/internals/usecase/session"
	"errors"
	"fmt"
	"strings"
)

type Content struct {
	Session          *session.SessionUseCase
	MemberAdapter    adapters.MemberAdapter
	CommunityAdapter adapters.CommunityAdapter
	PostAdapter      adapters.PostAdapter
}

func (this *Content) CreateCommunity(name, description, token string, nsfw bool) (*models.Community, error) {
	user, rewokeErr := this.Session.Rewoke(token)

	if rewokeErr != nil {
		return nil, rewokeErr
	}

	community, communityErr := models.NewCommunity(strings.ToLower(name), description, user.ID, nsfw)

	if communityErr != nil {
		return nil, errors.New("Couldn't create the community")
	}

	member, memberErr := models.NewMember(user.ID, community.ID)

	if memberErr != nil {
		return nil, errors.New("Couldn't create the community")
	}

	communityAdapterErr := this.CommunityAdapter.Add(community)

	if communityAdapterErr != nil {
		return nil, errors.New("Couldn't save the community")
	}

	memberAdapterErr := this.MemberAdapter.Add(member)

	if memberAdapterErr != nil {
		return nil, errors.New("COuldn't create the community")
	}

	return community, nil
}

//That functions is responsible to execute the operation to a user join in a community
func (this *Content) JoinTheCommunity(token string, communityName string) (*models.Member, error) {
	user, rewokeErr := this.Session.Rewoke(token)

	if rewokeErr != nil {
		return nil, rewokeErr
	}

	comm, commErr := this.CommunityAdapter.Name(strings.ToLower(communityName))

	if commErr != nil {
		return nil, errors.New("COuldn't find a community with this name")
	}

	_, memberAlreadyExists := this.MemberAdapter.MemberInTheCommunity(comm.ID, user.ID)

	if memberAlreadyExists == nil {
		return nil, errors.New("The user is already on the community")
	}

	member, memberErr := models.NewMember(user.ID, comm.ID)

	if memberErr != nil {
		return nil, errors.New("Couldn't join on the community")
	}

	adapterErr := this.MemberAdapter.Add(member)

	if adapterErr != nil {
		return nil, errors.New("Couldn't join on the community")
	}

	return member, nil
}

func (this *Content) GetCommunity(name string) (*models.Community, error) {
	community, cmmErr := this.CommunityAdapter.Name(strings.ToLower(name))

	if cmmErr != nil {
		return nil, errors.New("Couldn't find the community with this name")
	}

	return community, nil
}

func (this *Content) CreatePost(title, text, token, communityName string, nsfw bool) (*models.Post, error) {
	user, rewokeErr := this.Session.Rewoke(token)

	if rewokeErr != nil {
		return nil, rewokeErr
	}

	community, communityErr := this.CommunityAdapter.Name(strings.ToLower(communityName))

	if communityErr != nil {
		return nil, errors.New("Couldn't find the community with this name")
	}

	member, memberErr := this.MemberAdapter.MemberInTheCommunity(community.ID, user.ID)

	if memberErr != nil {
		return nil, errors.New("You don't have permission to create a post on this guild, since that you are not a member of it")
	}

	post, postErr := models.NewPost(title, text, member.ID, community.ID, nsfw)

	if postErr != nil {
		return nil, errors.New("Couldn't create the post")
	}

	adapterErr := this.PostAdapter.Add(post)

	if adapterErr != nil {
		return nil, adapterErr
	}

	return post, nil
}

func (this *Content) Search(text string) ([]models.Community, []models.Post, error) {
	communities, communitiesErr := this.CommunityAdapter.Search(text)

	if communitiesErr != nil {
		return nil, nil, errors.New("Couldn't find anything")
	}

	posts, postsErr := this.PostAdapter.Search(text)

	if postsErr != nil {
		return nil, nil, errors.New("Couldn't find anything")
	}

	return communities, posts, nil
}

func (this *Content) GetJoinedCommunities(token string) ([]models.Member, error) {
	user, rewokeErr := this.Session.Rewoke(token)

	if rewokeErr != nil {
		return nil, rewokeErr
	}

	members, memberErr := this.MemberAdapter.UserID(user.ID)

	fmt.Printf("%v\n", members)

	if memberErr != nil {
		return nil, errors.New("Couldn't find any community")
	}

	return members, nil
}
