package count

import (
	dg "github.com/bwmarrin/discordgo"
)

type EmojiCount struct {
	Emoji *dg.Emoji
	Count int64
}

func CountAllEmoji(s *dg.Session, chanID string, before string, emjs map[string]EmojiCount) (*map[string]EmojiCount, error) {
	msgs, err := s.ChannelMessages(chanID, 100, before, "", "")
	if err != nil {
		return nil, err
	}

	msgID := ""
	for _, msg := range msgs {
		for _, v := range msg.Reactions {
			if v.Emoji == nil {
				continue
			}

			if val, ok := emjs[v.Emoji.Name]; ok {
				val.Count += int64(v.Count)
				emjs[v.Emoji.Name] = val
			} else {
				emjs[v.Emoji.Name] = EmojiCount{
					Emoji: v.Emoji,
					Count: int64(v.Count),
				}
			}
		}
		msgID = msg.ID
	}
	if msgID != "" {
		CountAllEmoji(s, chanID, msgID, emjs)
	}

	res := sortEmj(&emjs)

	return res, nil
}

func sortEmj(emjs *map[string]EmojiCount) []*EmojiCount {
	res := make([]*EmojiCount, len(*emjs))
	for _, v := range *emjs {
		res = append(res, &v)
	}

	return res
}
