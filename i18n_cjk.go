package main

func init() {
	uiCopies["zh"] = uiCopy{
		MenuButton: "📋 菜单", StopButton: "⏹ 停止", Back: "◀ 返回", BackMenu: "◀ 返回菜单",
		MainMenuTitle: "主菜单", MainMenuBody: "请选择操作：",
		Learning: "学习", Words: "单词", Stats: "进度", Settings: "设置", Tools: "工具",
		NewLesson: "新课程", Practice: "练习", LevelTest: "水平", LearnWords: "学单词", WordGame: "复习", Spelling: "拼写",
		Vocabulary: "词汇表", Mistakes: "错误", Progress: "我的进度", Leaders: "排行榜", Limits: "限制",
		BotLanguage: "机器人语言", LearningLanguage: "学习语言", Notifications: "提醒", Premium: "Premium",
		ChooseBotLang: "选择机器人语言", ChooseLearnLang: "你想学习哪种语言？", ChooseTimezone: "选择时区",
		TimezoneHint: "我会按你的本地时间每天约 19:00 发送提醒。",
		BotLangSet:   "完成。机器人语言：%s。", LearnLangSet: "完成。现在学习：%s。",
		UnknownButton: "我没看懂这个按钮。请按 📋 菜单。", Stopped: "已停止。按 📋 菜单选择下一步。",
		WordQuestion: "%s 怎么说？", ChooseAnswer: "请选择正确答案：", Correct: "正确！", TryAgain: "还不对。再试一次：",
		NextWord: "下一个单词", NextPage: "下一页 ▶", AlreadyLearned: "这个单词已经在你的词汇表里。", AddedToVocab: "单词已加入已学词汇表。",
		TotalLearned: "已学单词：%d", WriteWord: "请写出 %s：", Hint: "提示", GoodSpelling: "拼写很好！",
	}
	uiCopies["ja"] = uiCopy{
		MenuButton: "📋 メニュー", StopButton: "⏹ 停止", Back: "◀ 戻る", BackMenu: "◀ メニューへ",
		MainMenuTitle: "メインメニュー", MainMenuBody: "操作を選んでください：",
		Learning: "学習", Words: "単語", Stats: "進捗", Settings: "設定", Tools: "ツール",
		NewLesson: "新しいレッスン", Practice: "練習", LevelTest: "レベル", LearnWords: "単語を学ぶ", WordGame: "復習", Spelling: "スペル",
		Vocabulary: "単語帳", Mistakes: "間違い", Progress: "進捗", Leaders: "ランキング", Limits: "制限",
		BotLanguage: "ボットの言語", LearningLanguage: "学習言語", Notifications: "通知", Premium: "Premium",
		ChooseBotLang: "ボットの言語を選んでください", ChooseLearnLang: "どの言語を学びたいですか？", ChooseTimezone: "タイムゾーンを選んでください",
		TimezoneHint: "あなたの時間で毎日19:00頃にリマインダーを送ります。",
		BotLangSet:   "完了しました。ボットの言語：%s。", LearnLangSet: "完了しました。これから学ぶ言語：%s。",
		UnknownButton: "そのボタンは理解できませんでした。📋 メニューを押してください。", Stopped: "停止しました。次の操作は 📋 メニューから選んでください。",
		WordQuestion: "%s はどう言いますか？", ChooseAnswer: "正しい答えを選んでください：", Correct: "正解！", TryAgain: "まだ違います。もう一度試してください：",
		NextWord: "次の単語", NextPage: "次へ ▶", AlreadyLearned: "この単語はすでに単語帳にあります。", AddedToVocab: "単語を学習済み単語帳に追加しました。",
		TotalLearned: "学習済み単語：%d", WriteWord: "%s を書いてください：", Hint: "ヒント", GoodSpelling: "正しく書けました！",
	}
	uiCopies["ko"] = uiCopy{
		MenuButton: "📋 메뉴", StopButton: "⏹ 중지", Back: "◀ 뒤로", BackMenu: "◀ 메뉴로",
		MainMenuTitle: "메인 메뉴", MainMenuBody: "작업을 선택하세요:",
		Learning: "학습", Words: "단어", Stats: "진도", Settings: "설정", Tools: "도구",
		NewLesson: "새 수업", Practice: "연습", LevelTest: "레벨", LearnWords: "단어 학습", WordGame: "복습", Spelling: "철자",
		Vocabulary: "단어장", Mistakes: "오답", Progress: "내 진도", Leaders: "순위", Limits: "한도",
		BotLanguage: "봇 언어", LearningLanguage: "학습 언어", Notifications: "알림", Premium: "Premium",
		ChooseBotLang: "봇 언어를 선택하세요", ChooseLearnLang: "어떤 언어를 배우고 싶나요?", ChooseTimezone: "시간대를 선택하세요",
		TimezoneHint: "현지 시간 기준 매일 약 19:00에 알림을 보냅니다.",
		BotLangSet:   "완료. 봇 언어: %s.", LearnLangSet: "완료. 이제 학습 언어: %s.",
		UnknownButton: "이 버튼을 이해하지 못했습니다. 📋 메뉴를 눌러 주세요.", Stopped: "중지되었습니다. 다음 작업은 📋 메뉴에서 선택하세요.",
		WordQuestion: "%s 는 어떻게 말하나요?", ChooseAnswer: "정답을 선택하세요:", Correct: "정답입니다!", TryAgain: "아직 아니에요. 다시 시도하세요:",
		NextWord: "다음 단어", NextPage: "다음 ▶", AlreadyLearned: "이 단어는 이미 단어장에 있습니다.", AddedToVocab: "단어가 학습한 단어장에 추가되었습니다.",
		TotalLearned: "학습한 단어: %d", WriteWord: "%s 를 쓰세요:", Hint: "힌트", GoodSpelling: "철자가 정확합니다!",
	}

	toolUICopies["zh"] = toolUICopy{
		VoiceToText: "语音转文字", ImageTranslate: "图片文字翻译", Translator: "翻译器", GPTAgent: "GPT 助手",
		VoicePrompt: "发送语音消息，我会把它转成文字。", ImagePrompt: "发送带文字的图片，我会识别并翻译。",
		TranslatorPrompt: "选择语言，然后发送文字、语音或带文字的照片。",
		VoiceModePrompt:  "发送语音消息。我会转成文字，不会启动练习。", ImageModePrompt: "发送带文字的图片。我会读取并翻译。",
		TranslatorModePrompt:   "翻译器已准备好。选择源语言和目标语言，然后发送文字、语音或图片。",
		ImageOpenToolsPrompt:   "要翻译图片中的文字，请打开 🧰 工具并选择图片文字翻译。",
		ImageSinglePhotoPrompt: "请只发送 1 张图片。打开 🧰 工具 -> 图片文字翻译，并发送一张照片。",
		TranscriptLabel:        "转写", TranslationLabel: "翻译",
		SourceLanguage: "源语言", TargetLanguage: "目标语言", AutoDetect: "自动", WebApp: "Web 应用",
		VoiceDiscussPrompt: "你可以用文字继续讨论，我会继续练习模式。", ImageDiscussPrompt: "你可以再发一张照片，我会处理新图片。也可以发文字，我会继续练习模式。",
		VoicePremiumRequired: "语音练习仅 Premium 可用。\n\nFree 目前只支持文字。\n打开 📋 菜单 -> Premium。",
		ImagePremiumRequired: "图片文字翻译仅 Premium 可用。\n\n打开 📋 菜单 -> Premium。",
		VoiceLimitReached:    "今天的语音消息额度已用完。\n\n语音消息：%d/%d\n明天会自动重置。",
		VoiceTooLong:         "语音消息太长。\n\n最大：%d 秒\n你的消息：%d 秒",
		VoiceDownloadFailed:  "无法下载语音消息：%s", VoiceTranscribeFailed: "无法识别语音消息：%s", VoiceTranslationFailed: "无法翻译转写文本：%s",
		ImageFileMissing: "找不到图片文件。请重新发送照片。", ImageDownloadFailed: "无法下载图片：%s", ImageReadFailed: "无法读取图片：%s",
	}
	toolUICopies["ja"] = toolUICopy{
		VoiceToText: "音声をテキスト化", ImageTranslate: "画像テキスト翻訳", Translator: "翻訳ツール", GPTAgent: "GPTエージェント",
		VoicePrompt: "音声メッセージを送ると、文字に起こします。", ImagePrompt: "文字入りの画像を送ると、読み取って翻訳します。",
		TranslatorPrompt: "言語を選び、テキスト、音声、または文字入り写真を送ってください。",
		VoiceModePrompt:  "音声メッセージを送ってください。文字起こしだけ行い、練習は開始しません。", ImageModePrompt: "文字入りの画像を送ってください。内容を読み取り翻訳します。",
		TranslatorModePrompt:   "翻訳ツールの準備ができました。元言語と翻訳先を選び、テキスト、音声、または画像を送ってください。",
		ImageOpenToolsPrompt:   "画像内の文字を翻訳するには、🧰 ツールを開いて画像テキスト翻訳を選んでください。",
		ImageSinglePhotoPrompt: "画像は1枚だけ送ってください。🧰 ツール -> 画像テキスト翻訳を開き、写真を1枚送ってください。",
		TranscriptLabel:        "文字起こし", TranslationLabel: "翻訳",
		SourceLanguage: "元の言語", TargetLanguage: "翻訳先の言語", AutoDetect: "自動", WebApp: "Webアプリ",
		VoiceDiscussPrompt: "この内容についてテキストで続けられます。練習モードは続きます。", ImageDiscussPrompt: "別の写真を送れば処理します。テキストで書けば練習モードを続けます。",
		VoicePremiumRequired: "音声練習は Premium で利用できます。\n\nFree は現在テキストのみ対応です。\n📋 メニュー -> Premium を開いてください。",
		ImagePremiumRequired: "画像テキスト翻訳は Premium で利用できます。\n\n📋 メニュー -> Premium を開いてください。",
		VoiceLimitReached:    "今日の音声メッセージ上限に達しました。\n\n音声メッセージ：%d/%d\n明日、自動でリセットされます。",
		VoiceTooLong:         "音声メッセージが長すぎます。\n\n最大：%d秒\nあなたのメッセージ：%d秒",
		VoiceDownloadFailed:  "音声メッセージをダウンロードできませんでした：%s", VoiceTranscribeFailed: "音声を認識できませんでした：%s", VoiceTranslationFailed: "文字起こしを翻訳できませんでした：%s",
		ImageFileMissing: "画像ファイルが見つかりません。写真をもう一度送ってください。", ImageDownloadFailed: "画像をダウンロードできませんでした：%s", ImageReadFailed: "画像を読み取れませんでした：%s",
	}
	toolUICopies["ko"] = toolUICopy{
		VoiceToText: "음성을 텍스트로", ImageTranslate: "이미지 텍스트 번역", Translator: "번역기", GPTAgent: "GPT 에이전트",
		VoicePrompt: "음성 메시지를 보내면 텍스트로 바꿔 드립니다.", ImagePrompt: "텍스트가 있는 이미지를 보내면 읽고 번역합니다.",
		TranslatorPrompt: "언어를 선택한 뒤 텍스트, 음성 또는 텍스트가 있는 사진을 보내세요.",
		VoiceModePrompt:  "음성 메시지를 보내세요. 텍스트로 변환하고 연습은 시작하지 않습니다.", ImageModePrompt: "텍스트가 있는 이미지를 보내세요. 내용을 읽고 번역합니다.",
		TranslatorModePrompt:   "번역기가 준비되었습니다. 원문 언어와 번역 언어를 선택한 뒤 텍스트, 음성 또는 이미지를 보내세요.",
		ImageOpenToolsPrompt:   "이미지의 텍스트를 번역하려면 🧰 도구를 열고 이미지 텍스트 번역을 선택하세요.",
		ImageSinglePhotoPrompt: "이미지는 1장만 보내 주세요. 🧰 도구 -> 이미지 텍스트 번역을 열고 사진 한 장을 보내세요.",
		TranscriptLabel:        "전사", TranslationLabel: "번역",
		SourceLanguage: "원문 언어", TargetLanguage: "번역 언어", AutoDetect: "자동", WebApp: "웹 앱",
		VoiceDiscussPrompt: "텍스트로 계속 이야기할 수 있으며 연습 모드는 계속됩니다.", ImageDiscussPrompt: "다른 사진을 보내면 새 이미지를 처리합니다. 텍스트를 쓰면 연습 모드를 계속합니다.",
		VoicePremiumRequired: "음성 연습은 Premium에서 사용할 수 있습니다.\n\nFree는 현재 텍스트만 지원합니다.\n📋 메뉴 -> Premium을 여세요.",
		ImagePremiumRequired: "이미지 텍스트 번역은 Premium에서 사용할 수 있습니다.\n\n📋 메뉴 -> Premium을 여세요.",
		VoiceLimitReached:    "오늘 음성 메시지 한도를 모두 사용했습니다.\n\n음성 메시지: %d/%d\n내일 자동으로 초기화됩니다.",
		VoiceTooLong:         "음성 메시지가 너무 깁니다.\n\n최대: %d초\n보낸 메시지: %d초",
		VoiceDownloadFailed:  "음성 메시지를 다운로드할 수 없습니다: %s", VoiceTranscribeFailed: "음성을 인식할 수 없습니다: %s", VoiceTranslationFailed: "전사 텍스트를 번역할 수 없습니다: %s",
		ImageFileMissing: "이미지 파일을 찾을 수 없습니다. 사진을 다시 보내 주세요.", ImageDownloadFailed: "이미지를 다운로드할 수 없습니다: %s", ImageReadFailed: "이미지를 읽을 수 없습니다: %s",
	}

	systemUICopyOverrides["zh"] = systemUICopy{
		LessonKind: "课程", PracticeKind: "练习", SpamWait: "连续消息太多。请等待 %d 秒后再试。",
		UnknownInterfaceLang: "我没有理解机器人语言。请用 /uilanguage 重新选择。", UnknownLearningLang: "我没有理解学习语言。请用 /language 重新选择。", UnknownTimezone: "我没有理解时区。请用 /timezone 重新选择。",
		LevelManualButton: "选择我的水平", LevelStartButton: "开始测试", LevelStartText: "水平测试：%s\n\n如果你已经知道自己的水平，可以手动选择。如果不确定，请完成一个简短测试：我会问 %d 个 A1 到 C2 的问题，然后设置学习水平，课程和单词会按此调整。",
		ManualLevelText: "请选择你的学习水平。\n\n如果不确定，建议返回并完成测试。", LevelQuestionText: "问题 %d/%d\n\n%s", LevelCurrentScore: "当前分数：%d", LevelDontKnow: "不知道",
		LevelManualSet: "完成。你的学习水平已设为 %s。\n\n我会在课程、练习和新单词中使用它。你随时可以用 /level 重新测试。", LevelTestSet: "完成。你的学习水平是 %s。\n\n我会在课程、练习和新单词中使用它。单词学习会从这个水平及以上开始选择。",
		LevelUnavailable: "%s 的单独水平测试尚未连接。\n\n现在可以手动选择水平，AI 会在课程和练习中使用它。", LevelOpenFailed: "无法打开测试问题。请重新开始测试。", LevelUnknownChoice: "没有理解所选水平。请用 /level 重新选择。", LevelUnknownAnswer: "没有理解测试答案。请重新开始水平测试。", LevelInactive: "这个测试已失效。请用 /level 重新开始。", LevelStale: "这个问题已经不是最新的。请回答最新的测试消息。",
		LessonAnswerInstruction: "请用一条消息回答。检查后课程会自动结束。", LessonDone: "课程完成。按 %s 继续下一步。", PracticeStarted: "练习模式已开启。请用 %s 写任意句子。我会简单回答，并用 %s 解释。要退出，请按 %s。", LimitReached: "今天的 %s 额度已用完。\n\n%s\n\n明天会自动重置。",
		ReminderTitle: "通知设置", ReminderStatus: "状态", ReminderEnabled: "已开启", ReminderDisabled: "已关闭", ReminderTimezone: "时区", ReminderDefault: "默认", ReminderHour: "发送时间", ReminderHint: "你可以开启或关闭每日提醒并选择时区。", ReminderOnButton: "开启通知", ReminderOffButton: "关闭通知", ReminderFooter: "打开 %s 并选择下一步。可用 /reminderoff 关闭提醒。",
		ReminderMessages:   []string{"今天迈出一小步，比明天的完美计划更有用。打开一节课吧。", "每天几分钟练习，也能保持语言节奏。", "复习一个单词，就能让词汇继续生长。"},
		VocabLevelQuestion: "请选择 %s 中对应的正确单词：%s",
	}
	systemUICopyOverrides["ja"] = systemUICopy{
		LessonKind: "レッスン", PracticeKind: "練習", SpamWait: "連続メッセージが多すぎます。%d秒待ってからもう一度試してください。",
		UnknownInterfaceLang: "ボットの言語を理解できませんでした。/uilanguage で選び直してください。", UnknownLearningLang: "学習言語を理解できませんでした。/language で選び直してください。", UnknownTimezone: "タイムゾーンを理解できませんでした。/timezone で選び直してください。",
		LevelManualButton: "自分のレベルを選ぶ", LevelStartButton: "テスト開始", LevelStartText: "レベルチェック：%s\n\n自分のレベルが分かっている場合は手動で選べます。不確かな場合は短いテストを受けてください。A1からC2まで%d問を出し、学習レベルを設定してレッスンと単語を調整します。",
		ManualLevelText: "学習レベルを選んでください。\n\n迷う場合は戻ってテストを受けるのがおすすめです。", LevelQuestionText: "質問 %d/%d\n\n%s", LevelCurrentScore: "現在のスコア：%d", LevelDontKnow: "わかりません",
		LevelManualSet: "完了しました。学習レベルを %s に設定しました。\n\nレッスン、練習、新しい単語でこのレベルを使います。/level でいつでも再テストできます。", LevelTestSet: "完了しました。あなたの学習レベルは %s です。\n\nレッスン、練習、新しい単語でこのレベルを使います。単語学習はこのレベル以上から選ばれます。",
		LevelUnavailable: "%s の個別レベルテストはまだ接続されていません。\n\n今は手動でレベルを選べます。AIはレッスンと練習でそのレベルを使います。", LevelOpenFailed: "テスト問題を開けませんでした。テストをもう一度開始してください。", LevelUnknownChoice: "選択されたレベルを理解できませんでした。/level で選び直してください。", LevelUnknownAnswer: "テストの回答を理解できませんでした。レベルチェックをもう一度開始してください。", LevelInactive: "このテストはもう有効ではありません。/level で再開してください。", LevelStale: "この質問は最新ではありません。最新のテストメッセージに答えてください。",
		LessonAnswerInstruction: "1つのメッセージで返信してください。確認後、レッスンは自動で終了します。", LessonDone: "レッスン完了です。次へ進むには %s を押してください。", PracticeStarted: "練習モードがオンです。%s で好きな文を書いてください。簡単に返答し、%s で説明します。終了するには %s を押してください。", LimitReached: "今日の %s の上限に達しました。\n\n%s\n\n明日、自動でリセットされます。",
		ReminderTitle: "通知設定", ReminderStatus: "状態", ReminderEnabled: "オン", ReminderDisabled: "オフ", ReminderTimezone: "タイムゾーン", ReminderDefault: "既定", ReminderHour: "送信時刻", ReminderHint: "毎日の通知をオン/オフにし、タイムゾーンを選べます。", ReminderOnButton: "通知をオンにする", ReminderOffButton: "通知をオフにする", ReminderFooter: "%s を開いて次のステップを選んでください。/reminderoff でリマインダーをオフにできます。",
		ReminderMessages:   []string{"今日の小さな一歩は、明日の完璧な計画より役に立ちます。レッスンを開きましょう。", "数分の練習でも、言語のリズムを保てます。", "単語を一つ復習すると、語彙はまた育ちます。"},
		VocabLevelQuestion: "%s の正しい単語を選んでください：%s",
	}
	systemUICopyOverrides["ko"] = systemUICopy{
		LessonKind: "수업", PracticeKind: "연습", SpamWait: "연속 메시지가 너무 많습니다. %d초 기다린 뒤 다시 시도하세요.",
		UnknownInterfaceLang: "봇 언어를 이해하지 못했습니다. /uilanguage 로 다시 선택하세요.", UnknownLearningLang: "학습 언어를 이해하지 못했습니다. /language 로 다시 선택하세요.", UnknownTimezone: "시간대를 이해하지 못했습니다. /timezone 으로 다시 선택하세요.",
		LevelManualButton: "내 레벨 선택", LevelStartButton: "테스트 시작", LevelStartText: "레벨 확인: %s\n\n이미 레벨을 알고 있다면 직접 선택할 수 있습니다. 확실하지 않다면 짧은 테스트를 진행하세요. A1부터 C2까지 %d문항을 묻고 학습 레벨을 설정해 수업과 단어를 맞춥니다.",
		ManualLevelText: "학습 레벨을 선택하세요.\n\n확실하지 않다면 돌아가서 테스트를 보는 것이 좋습니다.", LevelQuestionText: "질문 %d/%d\n\n%s", LevelCurrentScore: "현재 점수: %d", LevelDontKnow: "모르겠어요",
		LevelManualSet: "완료. 학습 레벨을 %s 로 설정했습니다.\n\n수업, 연습, 새 단어에 이 레벨을 사용합니다. 언제든 /level 로 다시 테스트할 수 있습니다.", LevelTestSet: "완료. 학습 레벨은 %s 입니다.\n\n수업, 연습, 새 단어에 이 레벨을 사용합니다. 단어 학습은 이 레벨 이상에서 단어를 고릅니다.",
		LevelUnavailable: "%s 의 별도 레벨 테스트는 아직 연결되지 않았습니다.\n\n지금은 레벨을 직접 선택할 수 있고, AI가 수업과 연습에서 그 레벨을 사용합니다.", LevelOpenFailed: "테스트 문제를 열 수 없습니다. 테스트를 다시 시작하세요.", LevelUnknownChoice: "선택한 레벨을 이해하지 못했습니다. /level 로 다시 선택하세요.", LevelUnknownAnswer: "테스트 답변을 이해하지 못했습니다. 레벨 확인을 다시 시작하세요.", LevelInactive: "이 테스트는 더 이상 활성 상태가 아닙니다. /level 로 다시 시작하세요.", LevelStale: "이 질문은 최신 질문이 아닙니다. 최신 테스트 메시지에 답하세요.",
		LessonAnswerInstruction: "메시지 하나로 답하세요. 확인 후 수업은 자동으로 끝납니다.", LessonDone: "수업 완료. 다음 단계로 가려면 %s 를 누르세요.", PracticeStarted: "연습 모드가 켜졌습니다. %s 로 아무 문장이나 써 보세요. 간단히 답하고 %s 로 설명합니다. 나가려면 %s 를 누르세요.", LimitReached: "오늘의 %s 한도를 모두 사용했습니다.\n\n%s\n\n내일 자동으로 초기화됩니다.",
		ReminderTitle: "알림 설정", ReminderStatus: "상태", ReminderEnabled: "켜짐", ReminderDisabled: "꺼짐", ReminderTimezone: "시간대", ReminderDefault: "기본값", ReminderHour: "보낼 시간", ReminderHint: "매일 알림을 켜거나 끄고 시간대를 선택할 수 있습니다.", ReminderOnButton: "알림 켜기", ReminderOffButton: "알림 끄기", ReminderFooter: "%s 를 열고 다음 단계를 선택하세요. /reminderoff 로 알림을 끌 수 있습니다.",
		ReminderMessages:   []string{"오늘의 작은 언어 한 걸음이 내일의 완벽한 계획보다 낫습니다. 수업을 열어 보세요.", "몇 분의 연습만으로도 언어 감각을 유지할 수 있습니다.", "단어 하나를 복습하면 어휘가 계속 자랍니다."},
		VocabLevelQuestion: "%s 에서 알맞은 단어를 선택하세요: %s",
	}

	premiumUICopyOverrides["zh"] = premiumUICopy{
		FreeVoiceUnavailable: "语音消息不可用", LessonsPerDay: "每天 %d 节课程", PracticeMessagesPerDay: "每天 %d 条练习消息", PremiumVoicesPerDay: "每天 %d 条语音消息，最长 %d 秒", VoiceTextAndTranslation: "语音转文字并翻译所说内容", ImageTextTranslation: "图片文字翻译", VoicePhotoContextPractice: "使用语音或照片上下文练习",
		MonthPrice: "30 天：%d ₽ 或 %d Stars", YearPrice: "一年：%d ₽ 或 %d Stars（约省 %d%%）", InviteFree: "邀请朋友可免费获得 7 天", InvoiceDescription: "Premium：每天 %d 节课程、%d 条练习消息、%d 条最长 %d 秒的语音消息、语音转文字和图片文字翻译。", PaymentUnrecognized: "已收到付款，但无法识别套餐。请打开 %s -> Premium。", PreCheckoutFailed: "无法确认购买。", ActivatedTitle: "Premium 已激活：%s", ActivatedAvailable: "现在可用：每天 %d 节课程、%d 条练习消息和 %d 条语音消息", ActiveUntil: "Premium 可用至 %s。", CheckLimits: "打开 %s -> %s 查看剩余额度。", YooKassaUnavailable: "YooKassa 付款尚未配置。请检查 .env。", YooKassaCreateFailed: "无法创建 YooKassa 付款，请稍后再试。", YooKassaPaymentText: "*%s* ⭐\n\n价格：*%d ₽*\n付款方式：*通过 YooKassa 的 SBP*\n\n付款后 Premium 会自动开启。", InviteCreateFailed: "无法创建邀请链接：%s", InviteLinkText: "用此链接邀请朋友：\n%s\n\n新用户首次通过它打开机器人后，你将获得 7 天 Premium。", ReferralShareText: "邀请朋友加入 Poliglot AI：\n%s\n\nPoliglot AI 是用于学习语言的 Telegram 机器人和 Web 应用：短课程、练习、单词、翻译器、语音和图片工具。\n\n朋友通过链接加入后，你将获得 7 天 Premium。", ReferralInviter: "有新学习者通过你的链接加入。Premium 已延长 %d 天，现在到 %s。", ReferralInvitee: "你通过邀请加入。你的朋友已经获得 7 天 Premium。",
		LimitsTitle: "我的额度", PlanLabel: "套餐", LessonsToday: "今日课程", PracticeToday: "今日练习", VoicesToday: "今日语音", MonthStarsButton: "30 天：%d Stars", MonthRubButton: "30 天：%d ₽ 通过 YooKassa", YearStarsButton: "一年：%d Stars", YearRubButton: "一年：%d ₽ 通过 YooKassa", InviteFriendButton: "邀请朋友", PaySBPButton: "通过 SBP 支付", PayStarsButton: "支付 %d Stars", MonthTitle: "Premium 30 天", YearTitle: "Premium 一年", MonthDays: "30 天", YearDays: "365 天",
	}
	premiumUICopyOverrides["ja"] = premiumUICopy{
		FreeVoiceUnavailable: "音声メッセージは利用できません", LessonsPerDay: "1日%dレッスン", PracticeMessagesPerDay: "1日%d件の練習メッセージ", PremiumVoicesPerDay: "1日%d件、最大%d秒の音声メッセージ", VoiceTextAndTranslation: "音声の文字起こしと発話内容の翻訳", ImageTextTranslation: "画像内テキストの翻訳", VoicePhotoContextPractice: "音声または写真の文脈を使った練習",
		MonthPrice: "30日：%d ₽ または %d Stars", YearPrice: "1年：%d ₽ または %d Stars（約%d%%割引）", InviteFree: "友達を招待して7日間無料", InvoiceDescription: "Premium：1日%dレッスン、%d件の練習メッセージ、音声メッセージ%d件（最大%d秒）、音声文字起こしと画像テキスト翻訳。", PaymentUnrecognized: "支払いを受け取りましたが、プランを認識できませんでした。%s -> Premium を開いてください。", PreCheckoutFailed: "購入を確認できませんでした。", ActivatedTitle: "Premium が %s 有効になりました", ActivatedAvailable: "現在利用可能：1日%dレッスン、%d件の練習メッセージ、%d件の音声メッセージ", ActiveUntil: "Premium は %s まで利用できます。", CheckLimits: "%s -> %s を開いて残りの上限を確認してください。", YooKassaUnavailable: "YooKassa決済はまだ設定されていません。.envを確認してください。", YooKassaCreateFailed: "YooKassa決済を作成できませんでした。少し後で試してください。", YooKassaPaymentText: "*%s* ⭐\n\n価格：*%d ₽*\n支払い方法：*YooKassa経由SBP*\n\n支払い後、Premiumは自動で有効になります。", InviteCreateFailed: "招待リンクを作成できませんでした：%s", InviteLinkText: "このリンクで友達を招待してください：\n%s\n\n新しいユーザーが初めてこのリンクからボットを開くと、あなたに7日間のPremiumが付与されます。", ReferralShareText: "友達を Poliglot AI に招待しましょう：\n%s\n\nPoliglot AI は、AIで言語を学ぶための Telegram ボット兼Webアプリです。短いレッスン、練習、単語、翻訳、音声、写真ツールを使えます。\n\n友達がリンクから参加すると、あなたに7日間のPremiumが付与されます。", ReferralInviter: "あなたのリンクから新しい学習者が参加しました。Premiumが%d日延長され、%s までになりました。", ReferralInvitee: "招待から参加しました。友達はすでに7日間のPremiumを受け取りました。",
		LimitsTitle: "私の上限", PlanLabel: "プラン", LessonsToday: "今日のレッスン", PracticeToday: "今日の練習", VoicesToday: "今日の音声", MonthStarsButton: "30日：%d Stars", MonthRubButton: "30日：%d ₽ YooKassa", YearStarsButton: "1年：%d Stars", YearRubButton: "1年：%d ₽ YooKassa", InviteFriendButton: "友達を招待", PaySBPButton: "SBPで支払う", PayStarsButton: "%d Starsで支払う", MonthTitle: "Premium 30日", YearTitle: "Premium 1年", MonthDays: "30日", YearDays: "365日",
	}
	premiumUICopyOverrides["ko"] = premiumUICopy{
		FreeVoiceUnavailable: "음성 메시지를 사용할 수 없음", LessonsPerDay: "하루 %d개 수업", PracticeMessagesPerDay: "하루 %d개 연습 메시지", PremiumVoicesPerDay: "하루 %d개 음성 메시지, 최대 %d초", VoiceTextAndTranslation: "음성을 텍스트로 변환하고 말한 내용 번역", ImageTextTranslation: "이미지 텍스트 번역", VoicePhotoContextPractice: "음성 또는 사진 문맥으로 연습",
		MonthPrice: "30일: %d ₽ 또는 %d Stars", YearPrice: "1년: %d ₽ 또는 %d Stars (약 %d%% 할인)", InviteFree: "친구를 초대하고 7일 무료 받기", InvoiceDescription: "Premium: 하루 %d개 수업, %d개 연습 메시지, 음성 메시지 %d개(최대 %d초), 음성 텍스트 변환 및 이미지 텍스트 번역.", PaymentUnrecognized: "결제는 받았지만 요금제를 인식하지 못했습니다. %s -> Premium을 여세요.", PreCheckoutFailed: "구매를 확인할 수 없습니다.", ActivatedTitle: "Premium 활성화: %s", ActivatedAvailable: "이제 사용 가능: 하루 %d개 수업, %d개 연습 메시지, %d개 음성 메시지", ActiveUntil: "Premium은 %s까지 사용할 수 있습니다.", CheckLimits: "%s -> %s 를 열어 남은 한도를 확인하세요.", YooKassaUnavailable: "YooKassa 결제가 아직 설정되지 않았습니다. .env를 확인하세요.", YooKassaCreateFailed: "YooKassa 결제를 만들 수 없습니다. 잠시 후 다시 시도하세요.", YooKassaPaymentText: "*%s* ⭐\n\n가격: *%d ₽*\n결제 방법: *YooKassa를 통한 SBP*\n\n결제 후 Premium이 자동으로 켜집니다.", InviteCreateFailed: "초대 링크를 만들 수 없습니다: %s", InviteLinkText: "이 링크로 친구를 초대하세요:\n%s\n\n새 사용자가 처음 이 링크로 봇을 열면 당신은 Premium 7일을 받습니다.", ReferralShareText: "친구를 Poliglot AI에 초대하세요:\n%s\n\nPoliglot AI는 AI로 언어를 배우는 Telegram 봇이자 웹 앱입니다. 짧은 수업, 연습, 단어, 번역기, 음성 및 사진 도구를 제공합니다.\n\n친구가 링크로 가입하면 당신은 Premium 7일을 받습니다.", ReferralInviter: "당신의 링크로 새 학습자가 참여했습니다. Premium이 %d일 연장되어 이제 %s까지입니다.", ReferralInvitee: "초대로 참여했습니다. 친구는 이미 Premium 7일을 받았습니다.",
		LimitsTitle: "내 한도", PlanLabel: "요금제", LessonsToday: "오늘 수업", PracticeToday: "오늘 연습", VoicesToday: "오늘 음성", MonthStarsButton: "30일: %d Stars", MonthRubButton: "30일: %d ₽ YooKassa", YearStarsButton: "1년: %d Stars", YearRubButton: "1년: %d ₽ YooKassa", InviteFriendButton: "친구 초대", PaySBPButton: "SBP로 결제", PayStarsButton: "%d Stars 결제", MonthTitle: "Premium 30일", YearTitle: "Premium 1년", MonthDays: "30일", YearDays: "365일",
	}
}
