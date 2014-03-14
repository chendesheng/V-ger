package language

import (
	"testing"
)

func TestDetectTChinese(t *testing.T) {
	lang1, lang2 := DetectLanguages(`
曆史的謊言

事是由人做的，話是由人說的。
千年的曆史底蘊所促成的不是人格的升華，而是一種沒落；謊言是由卑劣編織的，而良知卻在喪失。國家的曆史，從某方麵來說，就是群體的曆史，而群體卻由少數的偽善者操縱。

我們很慶幸，我們還有曆史，為此，我們而驕傲，而沾沾自喜。但所不幸的是，曆史也是有謊言的。我們沉迷於奢華的個人追求，我們享受於靡醉的向往，但我們卻在不知不覺中沉淪。大我似乎已變得渺小不堪，而小我卻在肆無忌憚。 真相，是不為所知的，要是都知道了，更談何什麼秘密。為了某種利益，作為個人，我們都曾用謊言來偽裝，但偽裝的結果隨著時間的推移，會變得模糊不清。

我們都知道，謊言說過一百遍之後，就會成為真理。曆史，在某一特定階段，是需要謊言的，無有謊言的曆史，談不上真正的曆史。崇高的理想和追求不外乎是美麗的欺騙，燦爛的光環始終在陽光的照射下變得灰暗，然而，我們為這灰暗而喝彩。

謊言是一切卑鄙者和軟弱者的行為，也是一切陰暗物所共有的特點。真實的曆史，隻是一種人為的點綴，虛名的花環總有凋零的黃昏。曆史也許是可以容忍的，就是它真的在說謊話；而最怕的，而最不能容忍的，卻是我們的麻木和愚鈍。當然，醜陋的事我們是不願意承認的，更不願檢討，就象人一樣，知道自己的大糞本為肮髒之物，所以，必須關起廁所的門。

我們崇尚於我們的燦爛文化，我們膜拜於我們光輝的曆史，在一種不可抗拒的力量的作用下，我們忽然發覺，我們已沒有了自己的思想。曆史，的確是殘酷的，而真理卻在退縮。

好聽的話，我們都百聽不厭其煩，冠冕堂皇的曆史，也足以讓我們驕傲一生，就是真的藏垢納汙，那也是我們不願提及的。我們自恃為聰明，但我們真的忽視了自己的愚笨，無有個性和思想的人，隻配做動物。就象那被耍的猴子，讓它跪它跪著，讓它頭頂著磚，它頂著；過後呢？還要用繩子拴著。

真話是很難說的，說了真話的人，隻能讓其他的人恥笑；真理的道路是艱辛的，在真理的麵前，謬誤耀武揚威。我們都不願意過早的清醒，我們都習慣於閉起雙眼，有那麼幾個敢睜眼的，卻被整個的社會所遺棄，所抹殺，直至過多的人投來冷眼和諷刺。

曆史是會裝腔作勢的，它不停的繞著彎子，在眾口鑠金之中，它變成了陽萎。 我們很難逃避這個社會對我們的影響和束縛，我們在曆史的熏陶之下，都聰明的領悟到了生存之道。為自己的謊言而問心無愧，為自己的卑劣行徑而得意洋洋，殊不知，曆史終究有它清醒的時候，當它真的清醒之時，我們會發現曆史也會象《皇帝的新裝》上的那個小孩子，大喊著：他身上什麼也沒有穿！

為此，我為那少數的睜開眼的人及追求真理的人而驕傲，更為那眾多的閉著眼的人及還未清醒的人而羞愧，而悲哀！
`)
	if lang1 != "zh" {
		t.Errorf("Expect 'cn' but %s", lang1)
	}

	if lang2 != "" {
		t.Errorf("Expect empty but %s", lang2)
	}
}

func TestDetectSChinese(t *testing.T) {
	lang1, lang2 := DetectLanguages(`
历史的谎言
	
	事是由人做的，话是由人说的。
	千年的历史底蕴所促成的不是人格的升华，而是一种没落；谎言是由卑劣编织的，而良知却在丧失。国家的历史，从某方面来说，就是群体的历史，而群体却由少数的伪善者操纵。
	
	我们很庆幸，我们还有历史，为此，我们而骄傲，而沾沾自喜。但所不幸的是，历史也是有谎言的。我们沉迷于奢华的个人追求，我们享受于靡醉的向往，但我们却在不知不觉中沉沦。大我似乎已变得渺小不堪，而小我却在肆无忌惮。 真相，是不为所知的，要是都知道了，更谈何什么秘密。为了某种利益，作为个人，我们都曾用谎言来伪装，但伪装的结果随着时间的推移，会变得模煳不清。
	
	我们都知道，谎言说过一百遍之后，就会成为真理。历史，在某一特定阶段，是需要谎言的，无有谎言的历史，谈不上真正的历史。崇高的理想和追求不外乎是美丽的欺骗，灿烂的光环始终在阳光的照射下变得灰暗，然而，我们为这灰暗而喝彩。
	
	谎言是一切卑鄙者和软弱者的行为，也是一切阴暗物所共有的特点。真实的历史，只是一种人为的点缀，虚名的花环总有凋零的黄昏。历史也许是可以容忍的，就是它真的在说谎话；而最怕的，而最不能容忍的，却是我们的麻木和愚钝。当然，丑陋的事我们是不愿意承认的，更不愿检讨，就象人一样，知道自己的大粪本为肮脏之物，所以，必须关起厕所的门。
	
	我们崇尚于我们的灿烂文化，我们膜拜于我们光辉的历史，在一种不可抗拒的力量的作用下，我们忽然发觉，我们已没有了自己的思想。历史，的确是残酷的，而真理却在退缩。
	
	好听的话，我们都百听不厌其烦，冠冕堂皇的历史，也足以让我们骄傲一生，就是真的藏垢纳污，那也是我们不愿提及的。我们自恃为聪明，但我们真的忽视了自己的愚笨，无有个性和思想的人，只配做动物。就象那被耍的猴子，让它跪它跪着，让它头顶着砖，它顶着；过后呢？还要用绳子拴着。
	
	真话是很难说的，说了真话的人，只能让其他的人耻笑；真理的道路是艰辛的，在真理的面前，谬误耀武扬威。我们都不愿意过早的清醒，我们都习惯于闭起双眼，有那么几个敢睁眼的，却被整个的社会所遗弃，所抹杀，直至过多的人投来冷眼和讽刺。
	
	历史是会装腔作势的，它不停的绕着弯子，在众口铄金之中，它变成了阳萎。 我们很难逃避这个社会对我们的影响和束缚，我们在历史的熏陶之下，都聪明的领悟到了生存之道。为自己的谎言而问心无愧，为自己的卑劣行径而得意洋洋，殊不知，历史终究有它清醒的时候，当它真的清醒之时，我们会发现历史也会象《皇帝的新装》上的那个小孩子，大喊着：他身上什么也没有穿！
	
	为此，我为那少数的睁开眼的人及追求真理的人而骄傲，更为那众多的闭着眼的人及还未清醒的人而羞愧，而悲哀！

It is done by the people , by the people say the words .
It is done by the people , by the people say the words .
It is done by the people , by the people say the words .
`)
	if lang1 != "zh" {
		t.Errorf("Expect 'cn' but %s", lang1)
	}

	if lang2 != "" {
		t.Errorf("Expect empty but %s", lang2)
	}
}

func TestDetectEnglish(t *testing.T) {
	lang1, lang2 := DetectLanguages(`
History Lies

It is done by the people , by the people say the words .
Millennium 's historical heritage have contributed to the sublimation of personality is not , but a decline ; lies woven by the despicable , but rather the conscience of the lost. The country's history , from a sense that history groups , and manipulated by a handful of groups but hypocrites .

We are very happy that we have a history , and we are proud , complacent . But that , unfortunately, history is also a lie . We indulged in the luxury of a personal pursuit, we enjoy in extravagant drunk yearning , but we sink imperceptibly . I seem to have become big small bear , and small I was in unbridled . The truth is not known , and if all know , more to talk about a secret. For some benefits , as individuals, we have to use lies to camouflage, but disguised as a result of the passage of time , will become die burnt unclear.

We all know that after the lies have said a hundred times , it will become the truth . Historically, at a particular stage , is the need to lie , there lies no history , no real history. Lofty ideals and the pursuit of beauty is nothing more than deception, brilliant aura always becomes dark in the sunlight , however, we applaud this gloomy .

Lie is despicable and weak all the behavior is common to all the characteristics of dark matter . Real history , but an artificial embellishment , vanity always have withered wreath at dusk . History may be tolerable , is it really a liar ; and fear , and the most intolerable , but it is our numb and dull . Of course, the ugly things we are willing to admit , and more reluctant to review , like everyone else, know their shit this is dirty thing , so you must shut the toilet door .

We advocate in our splendid culture , we worship our glorious history , the role of an irresistible force , we suddenly realize that we have no idea of ​​their own . History is indeed brutal , but the truth was in retreat.

Nice words , we are one hundred patiently , sounding history, but also enough to make us proud life is really filth , that is why we do not want to mention . We count on to be wise , but we really ignore their stupid , no personality and thoughtful people , with only do animals. Like that was playing monkey , let it kneeling kneel and let his head with bricks, it wore ; after it? Also tied with a rope.

The truth is hard to say that the truth of man , only to let other people ridiculed ; road of truth is difficult, in the face of truth , falsehood swagger . We are reluctant to awake early , we are accustomed to close their eyes , there are so few dare to open eyes , but was abandoned by the whole of society , has been denied until too many people voted to sit and satire.

History will be pretentious , it kept Rao Zhaowan child in Zhongkoushuojin being , it becomes impotent . Our society is difficult to escape the shackles of our influence and our history under the influence , are clever insight to survive. For his lie and have a clear conscience , for their despicable acts and elation , not knowing that history will eventually have its sober , and when it 's really clear , we will find that history will like "The Emperor's New Clothes " on the children , shouting : him nothing to wear !

To this end, I have a few of those who opened his eyes and pride in the pursuit of truth , the more people that many eyes closed and yet sober person and shame , and sorrow !

	为此，我为那少数的睁开眼的人及追求真理的人而骄傲，更为那众多的闭着眼的人及还未清醒的人而羞愧，而悲哀！
	为此，我为那少数的睁开眼的人及追求真理的人而骄傲，更为那众多的闭着眼的人及还未清醒的人而羞愧，而悲哀！
	为此，我为那少数的睁开眼的人及追求真理的人而骄傲，更为那众多的闭着眼的人及还未清醒的人而羞愧，而悲哀！
`)
	if lang1 != "en" {
		t.Errorf("Expect 'en' but %s", lang1)
	}

	if lang2 != "" {
		t.Errorf("Expect empty but %s", lang2)
	}
}

func TestDetectEnglishAndChinese(t *testing.T) {
	lang1, lang2 := DetectLanguages(`
History Lies

It is done by the people , by the people say the words .
Millennium 's historical heritage have contributed to the sublimation of personality is not , but a decline ; lies woven by the despicable , but rather the conscience of the lost. The country's history , from a sense that history groups , and manipulated by a handful of groups but hypocrites .

We are very happy that we have a history , and we are proud , complacent . But that , unfortunately, history is also a lie . We indulged in the luxury of a personal pursuit, we enjoy in extravagant drunk yearning , but we sink imperceptibly . I seem to have become big small bear , and small I was in unbridled . The truth is not known , and if all know , more to talk about a secret. For some benefits , as individuals, we have to use lies to camouflage, but disguised as a result of the passage of time , will become die burnt unclear.

We all know that after the lies have said a hundred times , it will become the truth . Historically, at a particular stage , is the need to lie , there lies no history , no real history. Lofty ideals and the pursuit of beauty is nothing more than deception, brilliant aura always becomes dark in the sunlight , however, we applaud this gloomy .

Lie is despicable and weak all the behavior is common to all the characteristics of dark matter . Real history , but an artificial embellishment , vanity always have withered wreath at dusk . History may be tolerable , is it really a liar ; and fear , and the most intolerable , but it is our numb and dull . Of course, the ugly things we are willing to admit , and more reluctant to review , like everyone else, know their shit this is dirty thing , so you must shut the toilet door .

We advocate in our splendid culture , we worship our glorious history , the role of an irresistible force , we suddenly realize that we have no idea of ​​their own . History is indeed brutal , but the truth was in retreat.

Nice words , we are one hundred patiently , sounding history, but also enough to make us proud life is really filth , that is why we do not want to mention . We count on to be wise , but we really ignore their stupid , no personality and thoughtful people , with only do animals. Like that was playing monkey , let it kneeling kneel and let his head with bricks, it wore ; after it? Also tied with a rope.

The truth is hard to say that the truth of man , only to let other people ridiculed ; road of truth is difficult, in the face of truth , falsehood swagger . We are reluctant to awake early , we are accustomed to close their eyes , there are so few dare to open eyes , but was abandoned by the whole of society , has been denied until too many people voted to sit and satire.

History will be pretentious , it kept Rao Zhaowan child in Zhongkoushuojin being , it becomes impotent . Our society is difficult to escape the shackles of our influence and our history under the influence , are clever insight to survive. For his lie and have a clear conscience , for their despicable acts and elation , not knowing that history will eventually have its sober , and when it 's really clear , we will find that history will like "The Emperor's New Clothes " on the children , shouting : him nothing to wear !

To this end, I have a few of those who opened his eyes and pride in the pursuit of truth , the more people that many eyes closed and yet sober person and shame , and sorrow !


历史的谎言
	
	事是由人做的，话是由人说的。
	千年的历史底蕴所促成的不是人格的升华，而是一种没落；谎言是由卑劣编织的，而良知却在丧失。国家的历史，从某方面来说，就是群体的历史，而群体却由少数的伪善者操纵。
	
	我们很庆幸，我们还有历史，为此，我们而骄傲，而沾沾自喜。但所不幸的是，历史也是有谎言的。我们沉迷于奢华的个人追求，我们享受于靡醉的向往，但我们却在不知不觉中沉沦。大我似乎已变得渺小不堪，而小我却在肆无忌惮。 真相，是不为所知的，要是都知道了，更谈何什么秘密。为了某种利益，作为个人，我们都曾用谎言来伪装，但伪装的结果随着时间的推移，会变得模煳不清。
	
	我们都知道，谎言说过一百遍之后，就会成为真理。历史，在某一特定阶段，是需要谎言的，无有谎言的历史，谈不上真正的历史。崇高的理想和追求不外乎是美丽的欺骗，灿烂的光环始终在阳光的照射下变得灰暗，然而，我们为这灰暗而喝彩。
	
	谎言是一切卑鄙者和软弱者的行为，也是一切阴暗物所共有的特点。真实的历史，只是一种人为的点缀，虚名的花环总有凋零的黄昏。历史也许是可以容忍的，就是它真的在说谎话；而最怕的，而最不能容忍的，却是我们的麻木和愚钝。当然，丑陋的事我们是不愿意承认的，更不愿检讨，就象人一样，知道自己的大粪本为肮脏之物，所以，必须关起厕所的门。
	
	我们崇尚于我们的灿烂文化，我们膜拜于我们光辉的历史，在一种不可抗拒的力量的作用下，我们忽然发觉，我们已没有了自己的思想。历史，的确是残酷的，而真理却在退缩。
	
	好听的话，我们都百听不厌其烦，冠冕堂皇的历史，也足以让我们骄傲一生，就是真的藏垢纳污，那也是我们不愿提及的。我们自恃为聪明，但我们真的忽视了自己的愚笨，无有个性和思想的人，只配做动物。就象那被耍的猴子，让它跪它跪着，让它头顶着砖，它顶着；过后呢？还要用绳子拴着。
	
	真话是很难说的，说了真话的人，只能让其他的人耻笑；真理的道路是艰辛的，在真理的面前，谬误耀武扬威。我们都不愿意过早的清醒，我们都习惯于闭起双眼，有那么几个敢睁眼的，却被整个的社会所遗弃，所抹杀，直至过多的人投来冷眼和讽刺。
	
	历史是会装腔作势的，它不停的绕着弯子，在众口铄金之中，它变成了阳萎。 我们很难逃避这个社会对我们的影响和束缚，我们在历史的熏陶之下，都聪明的领悟到了生存之道。为自己的谎言而问心无愧，为自己的卑劣行径而得意洋洋，殊不知，历史终究有它清醒的时候，当它真的清醒之时，我们会发现历史也会象《皇帝的新装》上的那个小孩子，大喊着：他身上什么也没有穿！
	
	为此，我为那少数的睁开眼的人及追求真理的人而骄傲，更为那众多的闭着眼的人及还未清醒的人而羞愧，而悲哀！

`)
	if lang1 != "en" {
		t.Errorf("Expect 'en' but %s", lang1)
	}

	if lang2 != "zh" {
		t.Errorf("Expect 'cn' but %s", lang2)
	}
}
