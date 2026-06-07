package main

import "strings"

type generatedLevelGap struct {
	Question string
	Options  []string
}

type generatedLevelMeaning struct {
	Question string
	Options  []string
}

type generatedLevelProfile struct {
	Sentences []string
	Gaps      []generatedLevelGap
	Meanings  []generatedLevelMeaning
}

func init() {
	for language, profile := range generatedLevelProfiles {
		levelAssessmentQuestionsByLanguage[language] = professionalAssessmentQuestions(generatedProfessionalAssessmentQuestions(profile))
	}
}

func generatedProfessionalAssessmentQuestions(profile generatedLevelProfile) []levelQuestion {
	questions := make([]levelQuestion, 0, 36)
	pointBands := []int{4, 5, 7, 9, 11, 13}
	for index := 0; index < 12 && index < len(profile.Sentences); index++ {
		questions = append(questions, levelQuestion{
			Question:     "Выбери правильное предложение:",
			Options:      generatedSentenceOptions(profile.Sentences[index], index),
			CorrectIndex: 0,
			Points:       pointBands[index/2],
		})
	}
	for index := 0; index < 12; index++ {
		gap := profile.Gaps[index%len(profile.Gaps)]
		questions = append(questions, levelQuestion{
			Question:     "Заполни пропуск: " + gap.Question,
			Options:      gap.Options,
			CorrectIndex: 0,
			Points:       pointBands[index/2],
		})
	}
	for index := 0; index < 6; index++ {
		sentence := profile.Sentences[(index+6)%len(profile.Sentences)]
		questions = append(questions, levelQuestion{
			Question:     "Выбери самый естественный вариант:",
			Options:      generatedSentenceOptions(sentence, index+12),
			CorrectIndex: 0,
			Points:       pointBands[index],
		})
	}
	for index := 0; index < 6; index++ {
		meaning := profile.Meanings[index%len(profile.Meanings)]
		questions = append(questions, levelQuestion{
			Question:     meaning.Question,
			Options:      meaning.Options,
			CorrectIndex: 0,
			Points:       pointBands[index],
		})
	}
	return questions
}

func generatedSentenceOptions(correct string, index int) []string {
	trimmed := strings.TrimSpace(correct)
	noDot := strings.TrimRight(trimmed, ".!?؟।")
	words := strings.Fields(noDot)
	wrongA := noDot
	wrongB := noDot
	wrongC := noDot
	if len(words) >= 3 {
		wrongA = strings.Join(words[1:], " ")
		wrongB = words[0] + " " + words[0] + " " + strings.Join(words[1:], " ")
		swapped := append([]string(nil), words...)
		swapped[0], swapped[1] = swapped[1], swapped[0]
		wrongC = strings.Join(swapped, " ")
	} else if len([]rune(noDot)) > 6 {
		runes := []rune(noDot)
		wrongA = string(runes[:len(runes)-1])
		wrongB = string(runes[1:])
		wrongC = string(runes[:len(runes)/2]) + " " + string(runes[len(runes)/2:])
	}
	options := []string{trimmed, wrongA, wrongB, wrongC}
	if index%2 == 1 {
		options[2], options[3] = options[3], options[2]
	}
	return options
}

var generatedLevelProfiles = map[string]generatedLevelProfile{
	"ar": {
		Sentences: []string{"أنا طالب.", "أتحدث العربية كل يوم.", "عندي وقت الآن.", "أشرب القهوة في الصباح.", "ذهبت إلى المدرسة أمس.", "هذا الكتاب مفيد.", "أريد أن أتعلم أكثر.", "إذا كان لدي وقت فسأقرأ.", "رغم التعب واصل العمل.", "يعتمد القرار على البيانات.", "تحتاج النتيجة إلى مراجعة إضافية.", "هذا التفسير دقيق ولكنه محدود."},
		Gaps: []generatedLevelGap{
			{"أنا ___ جديد.", []string{"طالب", "طالبةً", "طلاب", "طالبين"}},
			{"أذهب إلى العمل ___ الصباح.", []string{"في", "على", "من", "عن"}},
			{"هو ___ القهوة كل يوم.", []string{"يشرب", "شرب", "اشرب", "مشروب"}},
			{"هذا النص ___ من ذلك.", []string{"أوضح", "واضح", "الأوضح", "وضوح"}},
			{"إذا درست أكثر، ___ أفضل.", []string{"ستتحدث", "تحدثت", "تتحدث أمس", "تحدث"}},
			{"تم ___ الوثيقة أمس.", []string{"توقيع", "وقع", "يوقع", "توقيعًا"}},
		},
		Meanings: []generatedLevelMeaning{
			{"Что значит لكن؟", []string{"غير أن / مع ذلك", "لذلك", "دائما", "بدلا من ذلك"}},
			{"Что значит رغم ذلك؟", []string{"مع ذلك", "بسبب ذلك", "قبل ذلك", "فقط"}},
			{"Что значит يعتمد على؟", []string{"يتوقف على", "يرفض", "يترجم", "ينسى"}},
		},
	},
	"bn": {
		Sentences: []string{"আমি একজন ছাত্র।", "আমি প্রতিদিন বাংলা বলি।", "এখন আমার সময় আছে।", "আমি সকালে কফি পান করি।", "গতকাল আমি স্কুলে গিয়েছিলাম।", "এই বইটি উপকারী।", "আমি আরও শিখতে চাই।", "সময় থাকলে আমি পড়ব।", "ক্লান্ত হলেও সে কাজ চালিয়ে গেল।", "সিদ্ধান্তটি তথ্যের উপর নির্ভর করে।", "ফলাফলটি আরও পর্যালোচনা দরকার।", "এই ব্যাখ্যাটি সঠিক কিন্তু সীমিত।"},
		Gaps: []generatedLevelGap{
			{"আমি একজন ___।", []string{"ছাত্র", "ছাত্ররা", "ছাত্রকে", "ছাত্রের"}},
			{"আমি সকালে ___ পান করি।", []string{"কফি", "কফির", "কফিতে", "কফিকে"}},
			{"সে প্রতিদিন খবর ___।", []string{"পড়ে", "পড়েছিল", "পড়া", "পড়বে গতকাল"}},
			{"এই লেখা ওইটার চেয়ে ___।", []string{"আরও পরিষ্কার", "পরিষ্কারতম", "পরিষ্কারভাবে", "পরিষ্কাররা"}},
			{"যদি সময় থাকে, আমি তোমাকে ___।", []string{"ফোন করব", "ফোন করেছি", "ফোন করা", "ফোন করেছিলাম কাল"}},
			{"নথিটি গতকাল ___ হয়েছিল।", []string{"স্বাক্ষরিত", "স্বাক্ষর", "স্বাক্ষর করে", "স্বাক্ষর করবে"}},
		},
		Meanings: []generatedLevelMeaning{
			{"Что значит কিন্তু?", []string{"তবে / বিপরীতে", "তাই", "কখনও না", "এর বদলে"}},
			{"Что значит তবুও?", []string{"তারপরও", "কারণ", "আগে", "শুধু"}},
			{"Что значит নির্ভর করে?", []string{"ভিত্তি করে", "ত্যাগ করে", "অনুবাদ করে", "ভুলে যায়"}},
		},
	},
	"cs": {
		Sentences: []string{"Jsem student.", "Mluvím česky každý den.", "Teď mám čas.", "Ráno piju kávu.", "Včera jsem šel do školy.", "Tato kniha je užitečná.", "Chci se naučit víc.", "Kdybych měl čas, četl bych.", "I přes únavu pokračoval v práci.", "Rozhodnutí závisí na datech.", "Výsledek vyžaduje další kontrolu.", "Toto vysvětlení je přesné, ale omezené."},
		Gaps: []generatedLevelGap{
			{"Jsem ___.", []string{"student", "studenti", "studenta", "studentem jsem"}},
			{"Ráno ___ kávu.", []string{"piju", "pít", "pil jsem zítra", "piješ mě"}},
			{"Včera jsem ___ doma.", []string{"byl", "být", "budu", "jsem zítra"}},
			{"Tento text je ___ než tamten.", []string{"jasnější", "nejjasnější", "jasně", "jasný než"}},
			{"Kdybych měl čas, ___ víc.", []string{"četl bych", "čtu včera", "budu četl", "číst bych"}},
			{"Dokument byl včera ___.", []string{"podepsán", "podepsat", "podepisuje", "podpis"}},
		},
		Meanings: []generatedLevelMeaning{
			{"Что значит ale?", []string{"avšak / nicméně", "proto", "nikdy", "místo toho"}},
			{"Что значит přesto?", []string{"navzdory tomu", "kvůli tomu", "před tím", "pouze"}},
			{"Что значит záviset na?", []string{"být podmíněn něčím", "odmítnout", "přeložit", "zapomenout"}},
		},
	},
	"el": {
		Sentences: []string{"Είμαι φοιτητής.", "Μιλάω ελληνικά κάθε μέρα.", "Έχω χρόνο τώρα.", "Πίνω καφέ το πρωί.", "Χθες πήγα στο σχολείο.", "Αυτό το βιβλίο είναι χρήσιμο.", "Θέλω να μάθω περισσότερα.", "Αν είχα χρόνο, θα διάβαζα.", "Παρά την κούραση συνέχισε να εργάζεται.", "Η απόφαση εξαρτάται από τα δεδομένα.", "Το αποτέλεσμα χρειάζεται επιπλέον έλεγχο.", "Αυτή η εξήγηση είναι ακριβής αλλά περιορισμένη."},
		Gaps: []generatedLevelGap{
			{"Είμαι ___.", []string{"φοιτητής", "φοιτητές", "φοιτητή", "φοιτητής είμαι"}},
			{"Το πρωί ___ καφέ.", []string{"πίνω", "πίνει", "πιω", "ήπια αύριο"}},
			{"Χθες ___ στο σπίτι.", []string{"ήμουν", "είμαι", "θα είμαι", "να είμαι"}},
			{"Αυτό το κείμενο είναι ___ από εκείνο.", []string{"πιο σαφές", "σαφέστατο από", "σαφώς", "σαφής"}},
			{"Αν είχα χρόνο, ___ περισσότερο.", []string{"θα διάβαζα", "διαβάζω χθες", "θα διαβάσω χθες", "διάβασμα"}},
			{"Το έγγραφο ___ χθες.", []string{"υπογράφηκε", "υπογράφει", "υπογραφή", "υπογράψει"}},
		},
		Meanings: []generatedLevelMeaning{
			{"Что значит όμως?", []string{"ωστόσο / αλλά", "επομένως", "ποτέ", "αντί γι' αυτό"}},
			{"Что значит παρ' όλα αυτά?", []string{"παρά ταύτα", "εξαιτίας αυτού", "πριν από αυτό", "μόνο"}},
			{"Что значит εξαρτάται από?", []string{"βασίζεται σε", "αρνείται", "μεταφράζει", "ξεχνά"}},
		},
	},
	"hi": {
		Sentences: []string{"मैं छात्र हूँ।", "मैं हर दिन हिंदी बोलता हूँ।", "मेरे पास अभी समय है।", "मैं सुबह कॉफी पीता हूँ।", "कल मैं स्कूल गया था।", "यह किताब उपयोगी है।", "मैं और सीखना चाहता हूँ।", "अगर मेरे पास समय होता, तो मैं पढ़ता।", "थकान के बावजूद उसने काम जारी रखा।", "निर्णय डेटा पर निर्भर करता है।", "परिणाम को और समीक्षा चाहिए।", "यह व्याख्या सटीक है लेकिन सीमित है।"},
		Gaps: []generatedLevelGap{
			{"मैं ___ हूँ।", []string{"छात्र", "छात्रों", "छात्र को", "छात्र से"}},
			{"मैं सुबह कॉफी ___ हूँ।", []string{"पीता", "पीना", "पिया कल", "पीएगा"}},
			{"कल मैं घर पर ___।", []string{"था", "हूँ", "होना", "रहूँगा"}},
			{"यह पाठ उससे ___ है।", []string{"अधिक स्पष्ट", "सबसे स्पष्ट से", "स्पष्टता", "स्पष्टों"}},
			{"अगर समय होगा, मैं तुम्हें ___।", []string{"फोन करूँगा", "फोन किया", "फोन करना", "फोन करता कल"}},
			{"दस्तावेज़ पर कल ___ किए गए।", []string{"हस्ताक्षर", "हस्ताक्षरित करता", "हस्ताक्षर करना", "हस्ताक्षर होगा"}},
		},
		Meanings: []generatedLevelMeaning{
			{"Что значит लेकिन?", []string{"परंतु / मगर", "इसलिए", "कभी नहीं", "इसके बजाय"}},
			{"Что значит फिर भी?", []string{"इसके बावजूद", "इसके कारण", "इससे पहले", "केवल"}},
			{"Что значит निर्भर करता है?", []string{"आधारित होता है", "मना करता है", "अनुवाद करता है", "भूलता है"}},
		},
	},
	"hu": {
		Sentences: []string{"Diák vagyok.", "Minden nap beszélek magyarul.", "Most van időm.", "Reggel kávét iszom.", "Tegnap iskolába mentem.", "Ez a könyv hasznos.", "Többet szeretnék tanulni.", "Ha lenne időm, olvasnék.", "A fáradtság ellenére folytatta a munkát.", "A döntés az adatoktól függ.", "Az eredmény további ellenőrzést igényel.", "Ez a magyarázat pontos, de korlátozott."},
		Gaps: []generatedLevelGap{
			{"___ vagyok.", []string{"Diák", "Diákok", "Diákot", "Diákkal vagyok"}},
			{"Reggel kávét ___.", []string{"iszom", "inni", "ittam holnap", "iszik engem"}},
			{"Tegnap otthon ___.", []string{"voltam", "vagyok", "leszek", "lenni"}},
			{"Ez a szöveg ___, mint az.", []string{"érthetőbb", "legérthetőbb", "érthetően", "érthető mint"}},
			{"Ha lenne időm, többet ___.", []string{"olvasnék", "olvasok tegnap", "fogok olvastam", "olvasni lennék"}},
			{"A dokumentumot tegnap ___.", []string{"aláírták", "aláírni", "aláír", "aláírás"}},
		},
		Meanings: []generatedLevelMeaning{
			{"Что значит azonban?", []string{"viszont / ugyanakkor", "ezért", "soha", "ehelyett"}},
			{"Что значит ennek ellenére?", []string{"mindazonáltal", "emiatt", "ez előtt", "csak"}},
			{"Что значит függ valamitől?", []string{"attól függ", "elutasít", "lefordít", "elfelejt"}},
		},
	},
	"id": {
		Sentences: []string{"Saya seorang siswa.", "Saya berbicara bahasa Indonesia setiap hari.", "Saya punya waktu sekarang.", "Saya minum kopi pada pagi hari.", "Kemarin saya pergi ke sekolah.", "Buku ini berguna.", "Saya ingin belajar lebih banyak.", "Jika saya punya waktu, saya akan membaca.", "Meskipun lelah, dia terus bekerja.", "Keputusan itu bergantung pada data.", "Hasil itu perlu ditinjau lagi.", "Penjelasan ini tepat tetapi terbatas."},
		Gaps: []generatedLevelGap{
			{"Saya seorang ___.", []string{"siswa", "siswa-siswa", "siswa ke", "siswa dari"}},
			{"Saya ___ kopi pada pagi hari.", []string{"minum", "meminum kemarin", "minuman", "diminum"}},
			{"Kemarin saya ___ di rumah.", []string{"berada", "akan berada", "menjadi besok", "ada besok"}},
			{"Teks ini ___ daripada yang itu.", []string{"lebih jelas", "paling jelas daripada", "jelasnya", "jelas-jelas"}},
			{"Jika ada waktu, saya ___ kamu.", []string{"akan menelepon", "menelepon kemarin", "teleponkan", "ditelepon"}},
			{"Dokumen itu ___ kemarin.", []string{"ditandatangani", "menandatangani", "tanda tangan", "ditanda besok"}},
		},
		Meanings: []generatedLevelMeaning{
			{"Что значит tetapi?", []string{"namun", "karena itu", "tidak pernah", "sebagai gantinya"}},
			{"Что значит meskipun begitu?", []string{"walaupun demikian", "karena itu", "sebelumnya", "hanya"}},
			{"Что значит bergantung pada?", []string{"tergantung pada", "menolak", "menerjemahkan", "melupakan"}},
		},
	},
	"nl": {
		Sentences: []string{"Ik ben student.", "Ik spreek elke dag Nederlands.", "Ik heb nu tijd.", "Ik drink 's ochtends koffie.", "Gisteren ging ik naar school.", "Dit boek is nuttig.", "Ik wil meer leren.", "Als ik tijd had, zou ik lezen.", "Ondanks de vermoeidheid bleef hij werken.", "De beslissing hangt af van de gegevens.", "Het resultaat vereist extra controle.", "Deze uitleg is nauwkeurig maar beperkt."},
		Gaps: []generatedLevelGap{
			{"Ik ben ___.", []string{"student", "studenten", "student te", "student ben"}},
			{"Ik ___ 's ochtends koffie.", []string{"drink", "drinken", "dronk morgen", "gedronken ik"}},
			{"Gisteren ___ ik thuis.", []string{"was", "ben", "zal zijn", "zijn"}},
			{"Deze tekst is ___ dan die.", []string{"duidelijker", "duidelijkst dan", "duidelijk", "duidelijkheid"}},
			{"Als ik tijd had, ___ ik meer.", []string{"zou lezen", "lees gisteren", "zal gelezen", "lezen zou"}},
			{"Het document werd gisteren ___.", []string{"ondertekend", "ondertekenen", "ondertekent", "handtekening"}},
		},
		Meanings: []generatedLevelMeaning{
			{"Что значит maar?", []string{"echter", "daarom", "nooit", "in plaats daarvan"}},
			{"Что значит desondanks?", []string{"ondanks dat", "daardoor", "daarvoor", "alleen"}},
			{"Что значит afhangen van?", []string{"afhankelijk zijn van", "weigeren", "vertalen", "vergeten"}},
		},
	},
	"sv": {
		Sentences: []string{"Jag är student.", "Jag talar svenska varje dag.", "Jag har tid nu.", "Jag dricker kaffe på morgonen.", "I går gick jag till skolan.", "Den här boken är användbar.", "Jag vill lära mig mer.", "Om jag hade tid skulle jag läsa.", "Trots tröttheten fortsatte han att arbeta.", "Beslutet beror på data.", "Resultatet behöver granskas ytterligare.", "Den här förklaringen är exakt men begränsad."},
		Gaps: []generatedLevelGap{
			{"Jag är ___.", []string{"student", "studenter", "studenten mig", "student till"}},
			{"Jag ___ kaffe på morgonen.", []string{"dricker", "dricka", "drack i morgon", "druckit jag"}},
			{"I går ___ jag hemma.", []string{"var", "är", "ska vara", "vara"}},
			{"Den här texten är ___ än den där.", []string{"tydligare", "tydligast än", "tydligt", "tydlighet"}},
			{"Om jag hade tid, ___ jag mer.", []string{"skulle läsa", "läser i går", "ska läste", "läsa skulle"}},
			{"Dokumentet ___ i går.", []string{"undertecknades", "underteckna", "undertecknar", "underskrift"}},
		},
		Meanings: []generatedLevelMeaning{
			{"Что значит men?", []string{"dock", "därför", "aldrig", "i stället"}},
			{"Что значит trots det?", []string{"ändå", "på grund av det", "före det", "bara"}},
			{"Что значит bero på?", []string{"vara beroende av", "vägra", "översätta", "glömma"}},
		},
	},
	"ta": {
		Sentences: []string{"நான் ஒரு மாணவன்.", "நான் தினமும் தமிழ் பேசுகிறேன்.", "இப்போது எனக்கு நேரம் உள்ளது.", "காலையில் நான் காபி குடிக்கிறேன்.", "நேற்று நான் பள்ளிக்குச் சென்றேன்.", "இந்த புத்தகம் பயனுள்ளது.", "நான் இன்னும் கற்றுக்கொள்ள விரும்புகிறேன்.", "நேரம் இருந்தால் நான் படிப்பேன்.", "சோர்விருந்தாலும் அவன் வேலை தொடர்ந்தான்.", "முடிவு தரவுகளின் மீது சார்ந்துள்ளது.", "முடிவுக்கு மேலும் ஆய்வு தேவை.", "இந்த விளக்கம் துல்லியமானது ஆனால் வரம்புடையது."},
		Gaps: []generatedLevelGap{
			{"நான் ஒரு ___ .", []string{"மாணவன்", "மாணவர்கள்", "மாணவனை", "மாணவரால்"}},
			{"காலையில் நான் காபி ___ .", []string{"குடிக்கிறேன்", "குடித்தேன் நேற்று", "குடி", "குடிப்பது"}},
			{"நேற்று நான் வீட்டில் ___ .", []string{"இருந்தேன்", "இருக்கிறேன்", "இருப்பேன்", "இரு"}},
			{"இந்த உரை அதைவிட ___ .", []string{"தெளிவானது", "மிகத் தெளிவு விட", "தெளிவாக", "தெளிவுகள்"}},
			{"நேரம் இருந்தால் நான் மேலும் ___ .", []string{"படிப்பேன்", "படித்தேன் நேற்று", "படிக்க", "படிப்பு"}},
			{"ஆவணம் நேற்று ___ .", []string{"கையொப்பமிடப்பட்டது", "கையொப்பமிடு", "கையொப்பம்", "கையொப்பமிடும்"}},
		},
		Meanings: []generatedLevelMeaning{
			{"Что значит ஆனால்?", []string{"இருப்பினும்", "அதனால்", "ஒருபோதும் இல்லை", "அதற்கு பதில்"}},
			{"Что значит இருந்தாலும்?", []string{"அப்படியிருந்தும்", "அதனால்", "அதற்கு முன்", "மட்டும்"}},
			{"Что значит சார்ந்துள்ளது?", []string{"அடிப்படையாக உள்ளது", "மறுக்கிறது", "மொழிபெயர்க்கிறது", "மறக்கிறது"}},
		},
	},
	"te": {
		Sentences: []string{"నేను విద్యార్థిని.", "నేను ప్రతిరోజూ తెలుగు మాట్లాడుతాను.", "ఇప్పుడు నాకు సమయం ఉంది.", "ఉదయం నేను కాఫీ తాగుతాను.", "నిన్న నేను పాఠశాలకు వెళ్లాను.", "ఈ పుస్తకం ఉపయోగకరంగా ఉంది.", "నేను ఇంకా నేర్చుకోవాలనుకుంటున్నాను.", "సమయం ఉంటే నేను చదువుతాను.", "అలసట ఉన్నప్పటికీ అతను పని కొనసాగించాడు.", "నిర్ణయం డేటాపై ఆధారపడి ఉంటుంది.", "ఫలితానికి మరింత సమీక్ష అవసరం.", "ఈ వివరణ ఖచ్చితమైనది కానీ పరిమితమైనది."},
		Gaps: []generatedLevelGap{
			{"నేను ___ .", []string{"విద్యార్థిని", "విద్యార్థులు", "విద్యార్థిని కు", "విద్యార్థితో"}},
			{"ఉదయం నేను కాఫీ ___ .", []string{"తాగుతాను", "తాగాను నిన్న", "తాగు", "తాగడం"}},
			{"నిన్న నేను ఇంట్లో ___ .", []string{"ఉన్నాను", "ఉంటాను", "ఉండాలి", "ఉంటున్నాను రేపు"}},
			{"ఈ వాక్యం దానికంటే ___ .", []string{"స్పష్టంగా ఉంది", "స్పష్టమైనది కంటే", "స్పష్టత", "స్పష్టాలు"}},
			{"సమయం ఉంటే నేను ఇంకా ___ .", []string{"చదువుతాను", "చదివాను నిన్న", "చదవడం", "చదువు"}},
			{"పత్రం నిన్న ___ .", []string{"సంతకం చేయబడింది", "సంతకం చేయు", "సంతకం", "సంతకం చేస్తుంది"}},
		},
		Meanings: []generatedLevelMeaning{
			{"Что значит కానీ?", []string{"అయితే", "అందుకే", "ఎప్పుడూ కాదు", "దాని బదులు"}},
			{"Что значит అయినప్పటికీ?", []string{"అయినా సరే", "దాని కారణంగా", "దానికి ముందు", "మాత్రమే"}},
			{"Что значит ఆధారపడి ఉంటుంది?", []string{"దానిపై ఆధారపడుతుంది", "నిరాకరిస్తుంది", "అనువదిస్తుంది", "మరిచిపోతుంది"}},
		},
	},
	"th": {
		Sentences: []string{"ฉันเป็นนักเรียน.", "ฉันพูดภาษาไทยทุกวัน.", "ตอนนี้ฉันมีเวลา.", "ฉันดื่มกาแฟตอนเช้า.", "เมื่อวานฉันไปโรงเรียน.", "หนังสือเล่มนี้มีประโยชน์.", "ฉันอยากเรียนรู้มากขึ้น.", "ถ้าฉันมีเวลา ฉันจะอ่าน.", "แม้จะเหนื่อย เขาก็ทำงานต่อ.", "การตัดสินใจขึ้นอยู่กับข้อมูล.", "ผลลัพธ์ต้องตรวจสอบเพิ่มเติม.", "คำอธิบายนี้ถูกต้องแต่มีข้อจำกัด."},
		Gaps: []generatedLevelGap{
			{"ฉันเป็น ___ .", []string{"นักเรียน", "นักเรียนหลาย", "นักเรียนไป", "นักเรียนของ"}},
			{"ตอนเช้าฉัน ___ กาแฟ.", []string{"ดื่ม", "ดื่มแล้วพรุ่งนี้", "การดื่ม", "ถูกดื่ม"}},
			{"เมื่อวานฉัน ___ ที่บ้าน.", []string{"อยู่", "จะอยู่", "เป็นพรุ่งนี้", "อยู่แล้วพรุ่งนี้"}},
			{"ข้อความนี้ ___ กว่านั้น.", []string{"ชัดเจนกว่า", "ชัดเจนที่สุดกว่า", "อย่างชัดเจน", "ความชัดเจน"}},
			{"ถ้ามีเวลา ฉัน ___ มากขึ้น.", []string{"จะอ่าน", "อ่านเมื่อวาน", "การอ่าน", "ถูกอ่าน"}},
			{"เอกสารถูก ___ เมื่อวาน.", []string{"ลงนาม", "ลงนามพรุ่งนี้", "ลายเซ็น", "ลงนามอยู่"}},
		},
		Meanings: []generatedLevelMeaning{
			{"Что значит แต่?", []string{"อย่างไรก็ตาม", "ดังนั้น", "ไม่เคย", "แทนที่"}},
			{"Что значит ถึงอย่างนั้น?", []string{"อย่างไรก็ตาม", "เพราะสิ่งนั้น", "ก่อนหน้านั้น", "เท่านั้น"}},
			{"Что значит ขึ้นอยู่กับ?", []string{"ขึ้นกับ / พึ่งพา", "ปฏิเสธ", "แปล", "ลืม"}},
		},
	},
	"tl": {
		Sentences: []string{"Ako ay estudyante.", "Nagsasalita ako ng Tagalog araw-araw.", "May oras ako ngayon.", "Umiinom ako ng kape sa umaga.", "Pumunta ako sa paaralan kahapon.", "Kapaki-pakinabang ang aklat na ito.", "Gusto kong matuto pa.", "Kung may oras ako, magbabasa ako.", "Kahit pagod, nagpatuloy siya sa trabaho.", "Nakabatay sa datos ang desisyon.", "Kailangan pang suriin ang resulta.", "Tumpak ngunit limitado ang paliwanag na ito."},
		Gaps: []generatedLevelGap{
			{"Ako ay ___.", []string{"estudyante", "mga estudyante", "estudyante sa", "estudyante ng"}},
			{"___ ako ng kape sa umaga.", []string{"Umiinom", "Uminom bukas", "Inumin", "Iinom kahapon"}},
			{"Kahapon ___ ako sa bahay.", []string{"nasa", "ay bukas", "magiging", "maging"}},
			{"Mas ___ ang tekstong ito kaysa doon.", []string{"malinaw", "pinakamalinaw kaysa", "linaw", "malinaw na"}},
			{"Kung may oras ako, ___ ako.", []string{"magbabasa", "nagbasa bukas", "pagbasa", "binasa"}},
			{"Ang dokumento ay ___ kahapon.", []string{"nilagdaan", "lagda", "lalagdaan kahapon", "pumirma"}},
		},
		Meanings: []generatedLevelMeaning{
			{"Что значит ngunit?", []string{"subalit", "kaya", "hindi kailanman", "sa halip"}},
			{"Что значит gayunpaman?", []string{"sa kabila nito", "dahil dito", "bago nito", "lamang"}},
			{"Что значит nakabatay sa?", []string{"depende sa", "tumatanggi", "nagsasalin", "nakakalimot"}},
		},
	},
	"tr": {
		Sentences: []string{"Ben öğrenciyim.", "Her gün Türkçe konuşuyorum.", "Şu anda zamanım var.", "Sabah kahve içerim.", "Dün okula gittim.", "Bu kitap faydalı.", "Daha fazla öğrenmek istiyorum.", "Zamanım olsaydı okurdum.", "Yorgun olmasına rağmen çalışmaya devam etti.", "Karar verilere bağlıdır.", "Sonuç ek inceleme gerektiriyor.", "Bu açıklama doğru ama sınırlı."},
		Gaps: []generatedLevelGap{
			{"Ben ___.", []string{"öğrenciyim", "öğrenciler", "öğrenciye", "öğrenciyle"}},
			{"Sabah kahve ___.", []string{"içerim", "içmek", "içtim yarın", "içilmiş ben"}},
			{"Dün evde ___.", []string{"ydim", "yim", "olacağım dün", "olmak"}},
			{"Bu metin ondan daha ___.", []string{"açık", "en açık daha", "açıkça", "açıklık"}},
			{"Zamanım olsaydı daha çok ___.", []string{"okurdum", "okuyorum dün", "okuyacağım dün", "okumak"}},
			{"Belge dün ___.", []string{"imzalandı", "imzalamak", "imzalıyor", "imza"}},
		},
		Meanings: []generatedLevelMeaning{
			{"Что значит ancak?", []string{"fakat", "bu yüzden", "asla", "bunun yerine"}},
			{"Что значит buna rağmen?", []string{"yine de", "bu nedenle", "bundan önce", "sadece"}},
			{"Что значит bağlı olmak?", []string{"dayanmak", "reddetmek", "çevirmek", "unutmak"}},
		},
	},
	"vi": {
		Sentences: []string{"Tôi là học sinh.", "Tôi nói tiếng Việt mỗi ngày.", "Bây giờ tôi có thời gian.", "Tôi uống cà phê vào buổi sáng.", "Hôm qua tôi đã đi học.", "Cuốn sách này hữu ích.", "Tôi muốn học thêm.", "Nếu có thời gian, tôi sẽ đọc.", "Dù mệt, anh ấy vẫn tiếp tục làm việc.", "Quyết định phụ thuộc vào dữ liệu.", "Kết quả cần được xem xét thêm.", "Lời giải thích này chính xác nhưng có giới hạn."},
		Gaps: []generatedLevelGap{
			{"Tôi là ___.", []string{"học sinh", "các học sinh", "học sinh đến", "học sinh của"}},
			{"Buổi sáng tôi ___ cà phê.", []string{"uống", "uống hôm qua", "việc uống", "bị uống"}},
			{"Hôm qua tôi ___ ở nhà.", []string{"đã ở", "sẽ ở", "là ngày mai", "ở ngày mai"}},
			{"Văn bản này ___ văn bản kia.", []string{"rõ hơn", "rõ nhất hơn", "rõ ràng", "sự rõ"}},
			{"Nếu có thời gian, tôi ___ thêm.", []string{"sẽ đọc", "đã đọc ngày mai", "việc đọc", "bị đọc"}},
			{"Tài liệu đã được ___ hôm qua.", []string{"ký", "chữ ký", "ký ngày mai", "đang ký"}},
		},
		Meanings: []generatedLevelMeaning{
			{"Что значит nhưng?", []string{"tuy nhiên", "vì vậy", "không bao giờ", "thay vào đó"}},
			{"Что значит dù vậy?", []string{"tuy thế", "vì điều đó", "trước đó", "chỉ"}},
			{"Что значит phụ thuộc vào?", []string{"dựa vào / tùy thuộc vào", "từ chối", "dịch", "quên"}},
		},
	},
}
