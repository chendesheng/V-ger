package cld

import (
	"testing"
	"io/ioutil"
)
func TestDetectLanuage2(t *testing.T) {
	bytes, _ := ioutil.ReadFile("a.ass")
	text := string(bytes)
	lang1, lang2 := DetectLanguage2(text)
	println(lang1, lang2)
}
func Test2Langs(t *testing.T) {
	lang := DetectLanguage(`
1
0:01:05.00 --> 0:01:10.00
m 20 0 l 209 0 b 218 0 229 11 229 20 l 229 209 b 229 218 218 228 209 228 l 20 228 b 11 228 0 218 0 209 l 0 20 b 0 11 11 0 20 0

2
0:01:05.00 --> 0:01:10.00
m 115 64 l 115 35 l 171 0 l 177 0 l 133 46 l 143 52 l 179 0 l 184 0 l 159 63 l 229 107 l 229 182 l 145 96 l 135 122 l 207 227 l 137 227 l 115 172 l 115 93 l 126 76 l 115 64 l 104 76 l 115 93 l 115 172 l 93 227 l 23 227 l 95 122 l 84 96 l 0 182 l 0 107 l 70 63 l 45 0 l 50 0 l 87 52 l 97 46 l 52 0 l 59 0 l 115 35

3
0:01:05.00 --> 0:01:10.00
m 71 60 b 71 0 161 0 161 60 b 161 120 71 120 71 60

4
0:01:05.00 --> 0:01:10.00
m 325 0 l 355 0 l 332 47 l 354 47 l 399 136 l 369 136 l 328 53 l 288 136 l 257 136 l 325 0 m 467 0 l 497 0 l 474 47 l 496 47 l 541 136 l 511 136 l 470 53 l 430 136 l 399 136 l 467 0 m 545 1 l 583 1 l 583 14 l 568 14 l 568 19 l 583 19 l 583 30 l 568 30 l 568 34 l 599 34 l 599 30 l 583 30 l 583 19 l 599 19 l 599 14 l 583 14 l 583 1 l 611 1 b 616 1 622 6 622 10 l 622 36 l 652 0 l 678 0 l 644 41 l 622 41 l 622 47 l 596 47 l 597 54 l 625 54 l 625 68 l 541 68 l 541 54 l 572 54 l 571 47 l 545 47 l 545 1 m 583 72 l 583 85 l 569 85 l 569 90 l 598 90 l 598 85 l 583 85 l 583 72 l 611 72 b 615 72 621 78 621 82 l 653 44 l 678 44 l 644 86 l 621 86 l 621 103 l 597 103 l 597 136 l 570 136 l 564 126 l 562 136 l 542 136 l 548 107 l 568 107 l 565 121 l 571 121 l 571 103 l 547 103 l 547 72 l 583 72 m 600 107 l 620 107 l 624 124 l 653 89 l 679 89 l 642 136 l 615 136 l 618 132 l 606 132 l 600 107 m 689 0 l 716 0 l 721 15 l 732 15 l 732 30 l 718 56 l 731 56 l 735 100 l 721 100 l 717 59 l 714 64 l 714 136 l 693 136 l 693 79 l 676 79 l 707 30 l 679 30 l 679 15 l 694 15 l 689 0 m 738 0 l 804 0 b 807 0 813 6 813 9 l 813 87 l 794 87 l 794 14 l 756 14 l 756 87 l 763 77 l 763 21 l 787 21 l 787 91 l 798 91 l 798 120 l 820 120 l 812 136 l 778 136 l 778 90 l 748 136 l 723 136 l 756 87 l 738 87 l 738 0 m 257 151 l 275 151 l 297 182 l 319 151 l 337 151 l 304 197 l 304 227 l 290 227 l 290 197 l 257 151 m 337 151 l 355 151 l 377 182 l 399 151 l 417 151 l 384 197 l 384 227 l 370 227 l 370 197 l 337 151 m 425 192 l 447 192 b 445 181 427 181 425 192 l 410 192 b 414 162 458 162 462 192 l 462 206 l 425 206 b 425 210 429 212 433 213 l 462 213 l 462 227 l 433 227 b 412 227 408 203 410 192 m 462 151 l 532 151 l 532 165 l 504 165 l 504 227 l 489 227 l 489 165 l 462 165 l 462 151 m 580 172 l 580 186 l 549 186 b 543 186 543 192 549 192 l 565 192 b 589 192 589 227 565 227 l 532 227 l 532 213 l 565 213 b 570 213 570 206 565 206 l 549 206 b 524 206 524 172 549 172 l 580 172 m 592 213 l 606 213 l 606 227 l 592 227 l 592 213 m 639 172 l 665 172 l 665 186 l 639 186 b 623 186 623 213 639 213 l 665 213 l 665 227 l 639 227 b 603 227 603 172 639 172 m 700 184 b 722 184 722 215 700 215 l 700 229 b 740 229 740 170 700 170 b 660 170 660 229 700 229 l 700 215 b 680 215 680 184 700 184 m 737 172 l 782 172 b 803 172 813 177 813 198 l 813 228 l 799 228 l 799 198 b 799 186 793 187 782 187 l 782 228 l 768 228 l 768 187 l 752 187 l 752 228 l 737 228 l 737 172

5
0:01:05.00 --> 0:01:10.00
人人影视\NYYeTs.com

6
0:01:05.00 --> 0:01:10.00
■

7
0:01:05.00 --> 0:01:10.00
本字幕由人人影视字幕组翻译\N禁止用作任何商业盈利行为

8
0:01:10.50 --> 0:01:15.00
   翻译：冰蓝    Hannah    叉烧包    小花豆    北湾    喃喃

9
0:01:10.50 --> 0:01:15.00
■

10
0:01:15.50 --> 0:01:20.00
   时间轴：小萌    贝壳    后期：冰河   校对&监制：非燕

11
0:01:15.50 --> 0:01:20.00
■

12
0:00:25.00 --> 0:00:30.60
星  恋\N第一季 第一集

13
0:00:05.31 --> 0:00:09.80
2014年9月17日

14
0:00:13.44 --> 0:00:14.87
是我们抵达地球的日子\NThe day we arrived.

15
0:00:22.93 --> 0:00:24.65
我们从一个濒临灭亡的星球逃离\NWe fled a dying planet.

16
0:00:25.90 --> 0:00:27.65
我们的飞船迫降在这里\NOur ship crash-landed here.

17
0:00:32.49 --> 0:00:35.27
对于我们来说  这应是解放日\NFor my kind

18
0:00:38.14 --> 0:00:39.11
避难日\NRefuge.

19
0:00:47.41 --> 0:00:48.99
但是对于地球人来说\NBut for everyone on Earth...

20
0:00:50.59 --> 0:00:51.97
这却是侵略\Nit was an invasion.

21
0:00:57.90 --> 0:01:01.05
人类让我们别无选择  只能自我防卫\NThe humans left us no choice but to defend ourselves.

22
0:01:02.63 --> 0:01:03.98
弹药来袭\NIncoming!

23
0:01:23.04 --> 0:01:24.88
快跑  快跑

24
0:01:26.67 --> 0:01:27.53
快跑

25
0:01:32.50 --> 0:01:33.59
我听从了父亲的吩咐\NI did as my father said.

26
0:01:33.59 --> 0:01:35.91
我奋力飞奔\NI ran as fast as my legs would carry me

27
0:01:36.29 --> 0:01:38.50
不知道自己在逃避什么\Nnot understanding what I was running from.

28
0:01:39.60 --> 0:01:41.34
更不确定我在奔向什么\NAnd less sure what I was running to.

29
0:01:47.35 --> 0:01:48.59
国民警卫队继续\NNational Guardsmen continue

30
0:01:48.59 --> 0:01:51.85
和外星人进行一场流血冲突\Nto engage in a bloody conflict with the alien species

31
0:01:51.85 --> 0:01:55.07
他们的飞船坠落在巴吞鲁日外的几公里处\Nwhose ship crashed to Earth several miles from Baton Rouge.

32
0:01:55.07 --> 0:01:57.39
总统宣称进入紧急状态\NThe president has issued a state of emergency.

33
0:01:57.39 --> 0:01:59.68
而官方人士报告  一些侵略者\NWhile officials report that some of the invaders

34
0:01:59.68 --> 0:02:01.17
已经被包围\Nhave been rounded up

35
0:01:59.91 --> 0:02:02.76
侵略者袭击  伤亡人数增加

36
0:02:01.17 --> 0:02:03.83
可能有些外星人逃脱了追捕\Nit's possible some may have eluded capture.

37
0:02:03.83 --> 0:02:05.78
你觉得外星人的计划是什么\NWhat do you think would be the alien plan?

38
0:02:05.78 --> 0:02:07.03
我们该去我母亲家\NWe should go to my mother's.

39
0:02:07.03 --> 0:02:08.26
道路封锁了\NThe roads are closed.

40
0:02:08.26 --> 0:02:10.18
而且  去了那边情况也不会好转\NAnd plus

41
0:02:10.56 --> 0:02:12.06
尽量放松吧\NJust try and relax.

42
0:02:12.53 --> 0:02:13.65
怎么了  爸爸\NWhat's wrong

43
0:02:14.93 --> 0:02:16.02
没事  宝贝\NNothing

44
0:02:16.02 --> 0:02:17.31
一切都很好\NEverything's fine.

45
0:02:17.31 --> 0:02:18.82
别忘了吃药  埃默里\NDon't forget to take your medicine

46
0:02:20.26 --> 0:02:21.19
好的  爸爸\NOkay

47
0:02:24.35 --> 0:02:27.26
刚刚传来的消息  航空航天局的天文学家推断\NJust getting word; NASA astronomers have speculated

48
0:02:27.27 --> 0:02:29.82
这艘宇宙飞船的轨道\Nthat the trajectory of the spacecraft

49
0:02:29.82 --> 0:02:33.31
起源于一个名叫奥特里安的遥远的太阳系\Noriginated from a distant solar system called Atria.

50
0:02:33.64 --> 0:02:35.43
不管他们从何而来\NRegardless of where they came from

51
0:02:35.44 --> 0:02:39.02
当地居民和世界领导都很担心其意图\Nlocal residents and world leaders are wary of their intentions.

52
0:02:56.35 --> 0:02:57.31
有人吗\NHello?

53
0:02:59.94 --> 0:03:01.33
有人吗\NHello?

54
0:03:11.11 --> 0:03:12.87
没事的  没事的\NIt's okay. It's okay.

55
0:03:16.09 --> 0:03:17.65
你也是外星人  对吧\NYou're one of them

56
0:03:19.82 --> 0:03:21.14
我马上回来\NI'll be right back.

57
0:03:27.66 --> 0:03:28.47
另一方面  如果\NIf

58
0:03:28.47 --> 0:03:31.39
他们的目标实为征服并占领地球\Ntheir goal was actual occupation and conquest

59
0:03:31.39 --> 0:03:34.08
那他们可能会优先考虑那些...\Nthen they would probably have to prioritize anything...

60
0:03:36.13 --> 0:03:36.89
没事的\NIt's okay.

61
0:03:43.34 --> 0:03:45.00
你看起来不像个怪物\NYou don't look like a monster.

62
0:03:46.01 --> 0:03:47.33
你今晚可以待在这里\NYou can stay here tonight.

63
0:03:48.00 --> 0:03:49.25
我早晨再来看你\NI'll come see you in the morning.

64
0:03:49.25 --> 0:03:50.66
别出声\NBe very quiet.

65
0:04:19.27 --> 0:04:20.39
很好吃  对吧\NIt's good

66
0:04:21.92 --> 0:04:23.92
给  用叉子\NHere. Use the fork.

67
0:04:37.23 --> 0:04:38.52
这个叫星星\NThis one's called star.

68
0:04:38.93 --> 0:04:39.78
从这里穿过\NThrough here.

69
0:04:40.82 --> 0:04:42.29
从里面穿过\NThrough the loop.

70
0:04:44.07 --> 0:04:46.72
松开  就是颗星星  你来试试\NLet go

71
0:04:56.07 --> 0:04:57.89
B小队1-7\NBravo unit one-seven...

72
0:05:03.13 --> 0:05:04.28
得把你藏起来\NWe've got to hide you.

73
0:05:06.99 --> 0:05:07.99
他在这里\NHe's in here!

74
0:05:10.91 --> 0:05:12.77
快上  快上\NGo

75
0:05:12.77 --> 0:05:14.54
离她远点\NGet away from her!

76
0:05:14.54 --> 0:05:15.71
别伤害他\NDon't hurt him!

77
0:05:17.18 --> 0:05:19.19
-快  -不要\N- Now! - No!

78
0:06:04.95 --> 0:06:08.77
十年后

79
0:06:26.75 --> 0:06:28.82
十年前  名叫奥特里安人的\NIt's been ten years since the alien species

80
0:06:28.82 --> 0:06:30.82
外星物种抵达地球\Nknown as the Atrians arrived on Earth

81
0:06:30.82 --> 0:06:32.88
美国军方将其包围\Nand the U.S. Military rounded them up

82
0:06:32.88 --> 0:06:34.62
将他们送入一处政府辖区\Nand put them in a government facility

83
0:06:34.62 --> 0:06:36.40
一直关押到现在\Nwhere they remained in lockdown.

84
0:06:36.40 --> 0:06:39.87
今天  在尝试将他们融入社会的争议中\NToday

85
0:06:39.87 --> 0:06:42.83
七名奥特里安青少年会成为马歇尔高中的学生\Nseven Atrian teens will become students at Marshall High School.

86
0:06:40.00 --> 0:06:42.67
备受争议的取消七名奥特里安\N青少年隔离的计划

87
0:06:42.83 --> 0:06:43.70
小朱\NHey

88
0:06:44.22 --> 0:06:44.97
埃默里\NHey

89
0:06:45.96 --> 0:06:48.22
只有你才能在早上七点展现大汗淋漓的性感\NOnly you could pull off sweaty sexy at 7:00 A.M.

90
0:06:48.22 --> 0:06:51.15
在马歇尔高中外面  愤怒的抗议者要求\NOutside Marshall High

91
0:06:51.15 --> 0:06:54.10
奥特里安人继续被关押在他们的隔离区\Nthat the Atrians remain sequestered inside their sector.

92
0:06:50.95 --> 0:06:53.63
马歇尔高中实况\N愤怒的抗议仍在继续

93
0:06:54.10 --> 0:06:55.83
该做拉伸运动了  来吧\NTime to stretch. Come on.

94
0:06:55.83 --> 0:06:58.69
我的前任病友现在成了我的私人教练\NMy former hospital buddy's now my personal trainer.

95
0:06:59.13 --> 0:07:00.68
你说过一直想要一个的\NYou said you always wanted one.

96
0:07:00.68 --> 0:07:02.87
是啊  名字叫费尔南多\NYeah. Named Fernando.

97
0:07:02.87 --> 0:07:04.08
或者特里斯坦\NOr Tristan.

98
0:07:04.58 --> 0:07:06.40
他有外国口音  还有腹肌\NWith an accent and abs.

99
0:07:06.40 --> 0:07:09.10
这些反对消除隔离的拥护者相信奥特里安人\NThese anti-integration advocates believe the Atrians

100
0:07:08.08 --> 0:07:11.67
抗议者要求了解\N奥特里安人隐瞒了什么

101
0:07:09.10 --> 0:07:11.85
也许心怀鬼胎  想把地球变为殖民地\Nmay be hiding a greater plan to colonize Earth.

102
0:07:11.85 --> 0:07:14.54
如果要问我  这整个隔离计划\NIf you ask me

103
0:07:14.54 --> 0:07:16.95
只不过是为了分散大家的注意力\Nis just a distraction from the real issue.

104
0:07:16.95 --> 0:07:19.32
-每个频道都在放这个吗  -对\N- Is this on every channel? - Yep.

105
0:07:17.08 --> 0:07:19.46
市民对消除隔离的反应

106
0:07:20.91 --> 0:07:24.08
那艘飞船在整座小镇上投下了阴影\NThat ship is casting a shadow over this whole town...

107
0:07:21.20 --> 0:07:25.55
市民对消除隔离的反应

108
0:07:24.08 --> 0:07:25.49
疯子们今天都出来了\NAll the loonies are out today.

109
0:07:25.49 --> 0:07:27.69
真是不得不爱阴谋论者啊\NYou got to love conspiracy theorists.

110
0:07:27.69 --> 0:07:29.34
有一个人声称奥特里安人\NOne of them was claiming that the Atrians

111
0:07:29.34 --> 0:07:32.61
在隔离区里种植了一种草药  名叫柏草\Ngrow a medicinal herb inside their Sector called cyper.

112
0:07:33.91 --> 0:07:35.79
能治疗阴谋论者吗\NCan it cure conspiracy theorists?

113
0:07:35.79 --> 0:07:36.71
我只是说\NI'm just saying.

114
0:07:36.71 --> 0:07:39.14
也许你可以和一个奥特里安人搞好关系\NMaybe you could get chummy with one of the Atrians.

115
0:07:39.44 --> 0:07:40.41
拿到一批柏草\NGet a batch.

116
0:07:40.64 --> 0:07:41.88
拯救这整所医院\NFix up this whole hospital.

117
0:07:41.90 --> 0:07:42.79
我怎么说的\NWhat did I say?

118
0:07:43.04 --> 0:07:44.84
怀特希尔指挥官

119
0:07:43.44 --> 0:07:44.74
让他们后退\NKeep these guys back.

120
0:07:44.74 --> 0:07:45.76
天啊  是你爸爸\NOh

121
0:07:45.76 --> 0:07:49.26
如你所见  隔离区执法队的雷·怀特希尔指挥官\NAs you can see

122
0:07:46.25 --> 0:07:55.42
怀特希尔指挥官

123
0:07:49.27 --> 0:07:51.89
扫描了七名奥特里安人的腕带\Nscans the wristbands of the Atrian Seven

124
0:07:51.89 --> 0:07:55.53
这样他们在学校里时  当局就能追踪其位置\Nallowing authorities to track their location while at school.

125
0:07:55.53 --> 0:07:56.75
他真走运\NHe's so lucky.

126
0:07:56.75 --> 0:07:58.53
他是为数不多真正可以说\NHe's

127
0:07:58.53 --> 0:08:00.26
去过隔离区的人\Nthat they've been inside the Sector.

128
0:08:00.47 --> 0:08:02.96
当然只有你才会认为那是走运\NOf course only you would consider that lucky.

129
0:08:00.50 --> 0:08:03.26
七名奥特里安人离开隔离区

130
0:08:02.96 --> 0:08:05.60
那你对今天觉得紧张吗\NSo

131
0:08:06.95 --> 0:08:08.89
当然觉得\NYeah. Sure.

132
0:08:08.89 --> 0:08:10.70
昨晚新闻里的那个人说着\NI mean

133
0:08:10.70 --> 0:08:13.17
全世界都会关注马歇尔高中\Nhow the eyes of the world would be on Marshall High.

134
0:08:12.04 --> 0:08:14.21
七名奥特里安人离开隔离区

135
0:08:14.04 --> 0:08:16.73
让你四年后重返学校的第一天\NLeave it to you to start your first day back at school after four years

136
0:08:16.73 --> 0:08:18.11
就被全世界关注\Nin front of the eyes of the world.

137
0:08:18.11 --> 0:08:20.72
向政府施压\N...are putting pressure on our government.

138
0:08:20.72 --> 0:08:22.78
媒体没有进行全面报道\NThe media hasn't covered it thoroughly enough...

139
0:08:22.78 --> 0:08:24.21
真高兴你病情好转了  小埃\NI'm really glad you got better

140
0:08:26.88 --> 0:08:27.76
谢谢你\NThanks.

141
0:08:28.57 --> 0:08:29.56
轮到你了\NNow it's your turn.

142
0:08:31.03 --> 0:08:33.63
奥特里安七名青少年正前往马歇尔高中\NThe Atrian Seven are now headed to Marshall High

143
0:08:33.63 --> 0:08:35.77
那里的反对者越来越多\Nwhere protests continue to grow.

144
0:08:37.46 --> 0:08:39.11
我代表所有特区政府的同事们\NI speak for all my colleagues in D.C.

145
0:08:39.11 --> 0:08:41.77
通过将近三年的不断讨论\NWhen I say that after three years of ongoing discussion

146
0:08:41.77 --> 0:08:44.75
我们很荣幸  终于实现了这个项目\Nwe're very proud to see this program finally come to fruition.

147
0:08:44.91 --> 0:08:46.73
如果这一项目成功的话\NA program which

148
0:08:46.73 --> 0:08:50.13
可以让所有奥特里安人融入社会\Ncould lead to the integration of all Atrians into society.

149
0:08:50.13 --> 0:08:53.25
我们希望我们不仅能够和平共处\NThe hope is that we're not just able to peacefully coexist

150
0:08:53.25 --> 0:08:55.48
还能互相学习\Nbut we also learn from one another.

151
0:08:56.03 --> 0:08:57.41
他们的文化\NTheir culture alone...

152
0:08:57.41 --> 0:08:59.56
那些歪星人不属于这里\NThose Tatties belong with their own kind!

153
0:08:59.56 --> 0:09:01.24
这属于我们的孩子们\NOur kids were raised right.

154
0:09:01.24 --> 0:09:03.67
放开我  放开我\NGet off of me. Get off of me!

155
0:09:03.67 --> 0:09:04.73
歪星人  滚回家\NTatties

156
0:09:04.73 --> 0:09:06.41
走吧  咱们四处转转\NCome on

157
0:09:10.25 --> 0:09:12.50
科学实验室里有一些被编号的猴子\NIn the tech lab are code monkeys.

158
0:09:12.50 --> 0:09:14.16
就好像是这个奇怪的世界\NIt's kind of like this weird world

159
0:09:14.16 --> 0:09:16.39
被一些什么猴子啦猩猩啦统治了\Ndominated by chimps and orangutans or whatever--

160
0:09:16.39 --> 0:09:18.74
他们还不知道怎么和人交流\Nhas no idea how to communicate with humans.

161
0:09:18.74 --> 0:09:20.24
如果你看一眼\NAnd if you look discreetly

162
0:09:20.24 --> 0:09:22.63
你两点钟方向  亚洲时尚达人\Nto your two o'clock: Asian fashionistas.

163
0:09:22.63 --> 0:09:24.26
他们同样的衣服从来不穿两次\NThey never wear the same thing twice.

164
0:09:24.76 --> 0:09:26.42
看楼上  二楼平台上\NAnd further up the stairs

165
0:09:26.42 --> 0:09:28.04
是列宁崇拜者\Nare the Lenin worshipers.

166
0:09:28.66 --> 0:09:29.50
是那个共产主义列宁\NThe Communist Lenin

167
0:09:29.50 --> 0:09:30.33
不是披头士的那个\Nnot the Beatle.

168
0:09:30.76 --> 0:09:32.58
你第一节什么课  英语文学吗\NWhat's your first class? English lit?

169
0:09:32.58 --> 0:09:33.79
但愿你没选特克先生的课\NI hope you didn't get Mr. Turk.

170
0:09:33.80 --> 0:09:35.89
他会让你读《耸肩的阿特拉斯》\NI mean

171
0:09:35.89 --> 0:09:37.03
总之一句话\NAnd let's just say

172
0:09:37.03 --> 0:09:39.26
它的课特别沉闷\Nthat it is a total

173
0:09:39.26 --> 0:09:40.33
特别沉闷  小埃\NIt's a downer

174
0:09:40.33 --> 0:09:42.88
然后我肯定你没有听我说话\Nand I'm so sure that you're not even listening to me.

175
0:09:42.88 --> 0:09:45.65
-小埃  -我听着呢\N- Em! - Yes

176
0:09:45.65 --> 0:09:47.26
你刚刚说特别沉闷\NYou said it's a total and utter downer.

177
0:09:47.56 --> 0:09:48.41
他们来了\NThey're here.

178
0:10:03.33 --> 0:10:04.65
滚吧\NGet out of here!

179
0:10:04.82 --> 0:10:05.74
滚回家\NGo home!

180
0:10:05.74 --> 0:10:09.88
滚回家  滚回家\NGo home! Go home!

181
0:10:11.66 --> 0:10:12.55
都下车\NEverybody off.

182
0:10:13.95 --> 0:10:14.86
都下车\NEverybody off.

183
0:10:29.01 --> 0:10:30.00
好了  跟我来\NAll right

184
0:10:42.67 --> 0:10:43.97
来吧  快点\NLet's go. Move it.

185
0:10:43.97 --> 0:10:46.03
排好队  排成一排  肩并肩站好\NLine it up

186
0:10:46.49 --> 0:10:47.82
到这里排好\NI want a line right here.

187
0:10:49.01 --> 0:10:49.81
快排好\NLine it up.

188
0:10:50.40 --> 0:10:51.15
成一直线\NForm a line.

189
0:10:51.87 --> 0:10:54.18
在这儿排  你  站这里\NRight here along. You. Right here.

190
0:11:19.01 --> 0:11:21.75
今天还能再疯狂一点吗\NI mean

191
0:11:22.36 --> 0:11:23.59
刚才不就是吗\NI think it just did.

192
0:11:35.95 --> 0:11:37.18
我叫格洛丽亚\NMy name is Gloria.

193
0:11:37.59 --> 0:11:41.04
你们不认识我  但我熟悉你们每个人\NYou don't know me

194
0:11:41.23 --> 0:11:43.41
[外星语]

195
0:11:43.41 --> 0:11:45.30
没错  泰瑞  我最近是和人上床了\NYes

196
0:11:45.30 --> 0:11:48.86
并且在马歇尔高中只能说英语\NAnd English is the language of choice here at Marshall High.

197
0:11:50.35 --> 0:11:51.78
你们已经简单学习了礼仪\NI know you've been briefed on protocol

198
0:11:51.79 --> 0:11:53.76
通过了MHS入学考试\Nand passed your MHS entrance exams

199
0:11:53.76 --> 0:11:56.29
但是我认为有必要讲点规矩\Nbut I think it's necessary to set some ground rules

200
0:11:56.29 --> 0:11:57.53
为了你们的安全着想\Nfor your own safety.

201
0:11:57.53 --> 0:11:59.58
想必你们已经注意到增设的安全设施\NI'm sure you've noticed increased security.

202
0:11:59.59 --> 0:12:00.90
我们不希望有任何暴力事件\NWe don't anticipate any violence

203
0:12:00.90 --> 0:12:03.43
这些警卫是为了保护你们\Nbut the guards are here for your protection.

204
0:12:03.82 --> 0:12:04.97
和奥特里安人\NAnd the Atrians.

205
0:12:04.97 --> 0:12:06.83
警卫被批准使用震慑枪\NThe guards have been authorized to use their temblor guns

206
0:12:06.83 --> 0:12:08.19
如果你们走出警戒线\Nif you step out of line.

207
0:12:08.19 --> 0:12:10.89
你们只能离开隔离区去上学\NYou are only permitted to leave the Sector for school.

208
0:12:10.90 --> 0:12:12.31
宵禁仍然存在\NCurfew still stands.

209
0:12:12.31 --> 0:12:14.90
每晚必须在10点前回去\NBack in your pods by 10:00 p.m. every night.

210
0:12:14.90 --> 0:12:17.72
在学校里面  他们和你们一样\NBehind the walls of this school

211
0:12:17.72 --> 0:12:18.81
他们也是学生\NThey're students.

212
0:12:18.81 --> 0:12:20.33
你们是实验对象\NYou're test subjects.

213
0:12:20.33 --> 0:12:21.91
至少现在是\NAt least for the time being.

214
0:12:21.91 --> 0:12:23.98
但是总有一天  你们可以从这里毕业\NBut one day

215
0:12:24.18 --> 0:12:26.21
去上大学  或者工作\NHeading off to college or out in the work force.

216
0:12:26.21 --> 0:12:29.76
想一下  这个项目有可能成功\NImagine the possibilities if this program is a success.

217
0:12:29.76 --> 0:12:31.28
但是这很不容易\NBut it won't be easy.

218
0:12:31.28 --> 0:12:32.81
很多人都希望你们会失败\NA lot of people are hoping you fail.

219
0:12:32.81 --> 0:12:35.74
他们认为你们就是来毁灭我们的\NThey think your race has come here to destroy us.

220
0:12:35.74 --> 0:12:38.22
你们要证明他们说错了\NIt's up to you to prove them wrong.

221
0:12:47.05 --> 0:12:49.39
真希望咱们是那种拿着镭射枪的外星人\NWish we were the kind of aliens who carried ray guns.

222
0:12:49.40 --> 0:12:51.10
还记得咱们在隔离区学校的第一天吗\NRemember our first day at the Sector School?

223
0:12:51.10 --> 0:12:52.54
你们也不喜欢那里\NYou guys hated it there

224
0:12:52.54 --> 0:12:55.44
在那咱们至少不用担心震慑枪\NAt least we didn't have to worry about temblor guns pointed at our heads.

225
0:12:55.44 --> 0:12:57.82
或者更糟的  合唱团\NOr worse. Glee club.

226
0:13:00.42 --> 0:13:01.72
咱们的柜子还被做了标记\NThey tagged our lockers.

227
0:13:14.42 --> 0:13:15.18
右手\NRight.

228
0:13:15.76 --> 0:13:17.00
你得用右手按\NYou got to use your right hand.

229
0:13:20.03 --> 0:13:20.79
对了\NYes.

230
0:13:21.56 --> 0:13:24.14
我发誓我用高科技产品不经常这么差劲的\NI swear I'm not usually this technologically inept.

231
0:13:24.14 --> 0:13:24.76
这没什么\NIt's okay.

232
0:13:24.76 --> 0:13:26.28
你第一天表现得不错\NYou get a pass on your first day.

233
0:13:27.39 --> 0:13:28.69
你从哪转来的\NWhere you transferring from?

234
0:13:28.98 --> 0:13:30.57
我不是转学来的\NOh

235
0:13:30.57 --> 0:13:33.00
我只是过去四年\NI've just been out.

236
0:13:33.00 --> 0:13:34.65
一直在外面\NFor the last four years.

237
0:13:34.65 --> 0:13:37.26
我去年从纽霍尔转来的\NI transferred last year from Newhall.

238
0:13:37.26 --> 0:13:39.16
我也没参加迎新会什么的\NAnd I'm not on the welcoming committee or anything

239
0:13:39.16 --> 0:13:41.83
但我对这挺熟的  如果你需要什么帮忙\Nbut I know my way around

240
0:13:42.09 --> 0:13:42.91
谢谢你\NThanks.

241
0:13:43.58 --> 0:13:44.69
教我用手指\NFor the finger tip.

242
0:13:45.05 --> 0:13:48.34
我是指教我开柜子的方法\NI mean

243
0:13:49.09 --> 0:13:50.34
回见  再见\NSee you. Bye.

244
0:13:56.17 --> 0:13:57.56
我从来没想过我会这么说\NI never thought I'd say this

245
0:13:57.56 --> 0:13:59.67
但是我已经开始想念隔离区学校了\Nbut I'm already starting to miss our Sector School.

246
0:14:00.41 --> 0:14:01.68
人类就是野蛮人\NHumans are savages

247
0:14:01.68 --> 0:14:02.96
无一例外  罗曼\NRoman

248
0:14:04.43 --> 0:14:05.08
是吧\NRight?

249
0:14:06.61 --> 0:14:08.50
对  没错\NYeah. Right.

250
0:14:10.65 --> 0:14:13.07
我宣誓忠实于美利坚合众国国旗\NI pledge allegiance to the flag of the United of America

251
0:14:15.29 --> 0:14:18.00
忠实于她所代表的合众国\Nand to the republic

252
0:14:18.00 --> 0:14:20.28
苍天之下一个无可分割的国家\None nation

253
0:14:20.28 --> 0:14:22.56
在这里  人人享有自由与正义\Nwith liberty and justice for all.

254
0:14:25.25 --> 0:14:26.02
好了\NAll right.

255
0:14:26.02 --> 0:14:28.56
我将发给你们这个学期的课程表\NI'll be passing out your syllabus for the semester.

256
0:14:29.03 --> 0:14:30.25
自己留一份  然后往后传\NTake one and pass it back.

257
0:14:30.26 --> 0:14:32.62
据说他们有三个鸡鸡\NRumor is they have three penises.

258
0:14:32.62 --> 0:14:35.44
小号  中号  还有超大号\NSmall

259
0:14:37.37 --> 0:14:40.32
没错  她是泰勒  甜美  大姐大\NYep. Taylor. Sweet girl

260
0:14:40.32 --> 0:14:41.45
记得中学时代吧\NRemember in middle school?

261
0:14:41.69 --> 0:14:43.11
她现在掌权了\NYeah

262
0:14:43.64 --> 0:14:45.15
掌什么权\NIn charge of what?

263
0:14:45.15 --> 0:14:45.99
所有一切\NAll of it.

264
0:14:46.74 --> 0:14:48.49
你听说过奥特里安人\NHey

265
0:14:48.49 --> 0:14:51.03
可能在隔离区种了一种药草吗\Nthe Atrians supposedly grow in the Sector?

266
0:14:51.03 --> 0:14:53.31
好像叫"柏草"\NSomething called "Cyper"?

267
0:14:53.31 --> 0:14:55.14
据说这种药草有魔力\NIt's

268
0:14:55.14 --> 0:14:56.85
柏草  没听说过\NCyper? No.

269
0:14:56.85 --> 0:14:58.40
不过你要是指大麻\NBut if it's herbs you seek

270
0:14:58.40 --> 0:15:00.47
说不定瘾君子的柜子里有\Nsome of the stoners probably have some in their lockers.

271
0:15:04.48 --> 0:15:07.48
你们吃我们的食物  穿我们的衣服\NYou eat our food

272
0:15:07.48 --> 0:15:10.50
呼吸我们的空气  却不愿对我们的国旗宣誓\Nbreathe our air

273
0:15:10.50 --> 0:15:11.52
我们不是地球公民\NWe're not citizens.

274
0:15:11.52 --> 0:15:14.19
没错  你们是火星人\NOh

275
0:15:14.19 --> 0:15:17.39
实际上  火星人是火星来的\NWell

276
0:15:18.09 --> 0:15:19.66
我们的星球要高级多了\NWe're from a far superior planet.

277
0:15:19.67 --> 0:15:20.88
你说什么  怪胎\NWhat did you say to me

278
0:15:20.88 --> 0:15:22.99
到此为止\NHey

279
0:15:23.30 --> 0:15:25.87
奥特里安人不需要对国旗宣誓\NAtrians are not required to pledge to the flag.

280
0:15:26.64 --> 0:15:27.42
振作点\NCheer up.

281
0:15:27.63 --> 0:15:30.12
你们的自由和正义  我们不会夺去\NYou get to keep liberty and justice all to yourself.

282
0:15:30.12 --> 0:15:32.09
打开秋季学期日程表\NOpen up the fall semester calendar.

283
0:15:32.59 --> 0:15:34.06
有些重要活动\NWe have some events coming up.

284
0:15:35.33 --> 0:15:37.28
诺克斯  你一直极力主张奥特里安人\NNox

285
0:15:37.29 --> 0:15:38.44
与地球人享有平等权利\Nfor equal rights.

286
0:15:38.44 --> 0:15:40.71
你真的认为会有这么一天\NDid you think this day would ever come?

287
0:15:40.71 --> 0:15:41.98
我一直充满希望\NI never lost that hope.

288
0:15:42.24 --> 0:15:46.32
这个项目是迈向种族统一的第一步\NThis program is the first step toward unifying our races.

289
0:15:46.33 --> 0:15:48.23
我无比自豪\NAnd I couldn't be prouder

290
0:15:48.23 --> 0:15:51.81
因为我自己的孩子参与了这历史性的一刻\Nthat my own children are part of history in the making.

291
0:15:51.82 --> 0:15:53.57
-谢谢  先生  -谢谢你\N- Thank you

292
0:15:52.50 --> 0:15:57.46
奥特里安隔离区 - 东门

293
0:15:54.07 --> 0:15:55.47
这是奥特里安七人组  我们过去\NIt's the Atrian Seven! Let's go!

294
0:16:03.06 --> 0:16:04.37
感觉怎样\NHey. How did it go?

295
0:16:04.84 --> 0:16:07.37
很好  跟我们在隔离区受到的凌辱相比\NExcellent. We've been far too insulated in the Sector

296
0:16:07.37 --> 0:16:09.84
这里的人们善良多了\Nfrom just how kind and compassionate these people can be.

297
0:16:09.84 --> 0:16:11.23
你哥现在已经\NYour brother has mastered

298
0:16:11.23 --> 0:16:13.78
深谙地球人的幽默艺术了\Nthe earthly art of sarcasm with relative ease.

299
0:16:13.78 --> 0:16:16.15
说来听听  如果我有幸交到人类朋友\NTell me about it. If I have any chance of making friends

300
0:16:16.15 --> 0:16:17.89
我就装作不认识他了\NI'm gonna have to distance myself from him.

301
0:16:17.89 --> 0:16:19.38
为什么要跟人类做朋友\NWhy do you want to be friends with them?

302
0:16:19.38 --> 0:16:21.57
人类比咱们有趣多了\NHumans are so much more colorful.

303
0:16:21.57 --> 0:16:23.66
抱歉  但今天哪里有趣了\NI'm sorry. What about today was colorful?

304
0:16:23.67 --> 0:16:25.72
是武装警卫还是抗议者\NWas it the armed guards or the protestors?

305
0:16:25.73 --> 0:16:28.41
还是  我知道了  是憎恨我们的同学\NOr...? Oh

306
0:16:28.41 --> 0:16:29.35
他们只是还不了解我们\NThey just don't know us yet.

307
0:16:29.35 --> 0:16:30.71
[外星语]

308
0:16:30.71 --> 0:16:32.74
玛雅  讲英文\NMaia

309
0:16:33.98 --> 0:16:36.39
你们是去那里学习  不是交朋友的\NYou're there to learn

310
0:16:38.76 --> 0:16:40.43
明天又是新的一天\NWell

311
0:16:58.81 --> 0:17:00.52
你女儿自己会解释\NYour daughter will tell you exactly why

312
0:17:00.52 --> 0:17:03.09
她为什么咽不下虾球捞面\Nshe's not devouring that shrimp lo mein.

313
0:17:05.83 --> 0:17:07.03
没什么\NIt's nothing.

314
0:17:10.24 --> 0:17:15.00
我今天一直在想那个躲在棚屋里的男孩\NI just... was thinking a lot about that boy in the shed today.

315
0:17:15.72 --> 0:17:16.61
想他什么\NWhat about him?

316
0:17:17.34 --> 0:17:20.18
想如果他没有要保护我\NAbout how

317
0:17:20.18 --> 0:17:21.69
现在应该还活着\Nhe'd still be alive.

318
0:17:26.51 --> 0:17:28.58
不是所有人都能获救  宝贝\NNot everyone can be saved

319
0:17:28.58 --> 0:17:31.91
有时  要顾全大局\NSometimes... there's a greater plan at work.

320
0:17:36.05 --> 0:17:37.29
用餐愉快\NEnjoy your meal.

321
0:17:40.76 --> 0:17:42.65
下午好  奥特里安人\NGood afternoon

322
0:17:42.65 --> 0:17:45.36
你的午餐选择为爆米花蛋糕\NYour lunch option is puffed rice cake.

323
0:17:45.65 --> 0:17:48.97
你要鸡肉味  牛肉味  还是胡萝卜味\NWould you like chicken-flavored

324
0:17:50.64 --> 0:17:51.87
用餐愉快\NEnjoy your meal.

325
0:18:00.18 --> 0:18:01.94
你叫埃默里  对吧\NIt's Emery

326
0:18:01.94 --> 0:18:03.20
你住在伍德格伦\NYou live over on Woodglen?

327
0:18:03.20 --> 0:18:05.05
-是的  -手机拿出来\N- Yeah. - Let me see your phone.

328
0:18:08.14 --> 0:18:09.99
我们要在这个废弃农舍开派对\NThere's a party at this abandoned farmhouse.

329
0:18:09.99 --> 0:18:10.91
你也来吧\NYou should come.

330
0:18:12.26 --> 0:18:13.71
我可以带卢卡斯吗\NCan I bring Lukas?

331
0:18:13.71 --> 0:18:17.01
当然  如果你们俩是一对的话\NSure

332
0:18:17.01 --> 0:18:19.65
-我们是朋友  -好\N- We're friends. - Oh

333
0:18:20.21 --> 0:18:21.85
因为格雷森不确定\NBecause Grayson wasn't sure.

334
0:18:23.63 --> 0:18:26.97
格雷森要约我吗\NGrayson... was... asking?

335
0:18:35.00 --> 0:18:36.88
我来跟你换班\NI am relieving you of duties.

336
0:18:36.88 --> 0:18:39.97
好  你来得刚好\NYes

337
0:18:40.60 --> 0:18:42.30
那个大姐大要做什么\NSo

338
0:18:44.00 --> 0:18:45.56
你怎么会受到邀请  小埃\NHow did you manage this

339
0:18:45.88 --> 0:18:47.98
是她邀请我的  我们俩\NShe invited me. Or us.

340
0:18:48.31 --> 0:18:50.61
她肯定爱上你了  卢克\NI think she's got a crush on you

341
0:18:50.61 --> 0:18:52.26
不错啊  小埃  你真会哄人开心\NThat's cute

342
0:18:52.27 --> 0:18:53.86
这是你的才能吧  很好\NWas that your material? Great.

343
0:18:57.31 --> 0:18:58.71
好\NOkay.

344
0:19:01.81 --> 0:19:02.79
这是你的社团\NThis is your club?

345
0:19:05.01 --> 0:19:05.99
是的\NYeah.

346
0:19:05.99 --> 0:19:08.06
就是去本地的医院\NYou go to local hospitals

347
0:19:08.06 --> 0:19:11.06
陪病人画画  做剪贴\Nspend time with patients

348
0:19:11.06 --> 0:19:12.16
诸如此类\Nthat sort of thing.

349
0:19:16.19 --> 0:19:17.86
我觉得很抚慰人心\NI found it really comforting.

350
0:19:19.20 --> 0:19:20.53
你去年是社团成员\NYou were in the club last year?

351
0:19:21.06 --> 0:19:23.77
不是  我当时是病人\NNo. I was a patient.

352
0:19:24.30 --> 0:19:26.11
过去四年由于免疫缺陷\NI spent the last four years in the hospital

353
0:19:26.11 --> 0:19:27.76
我一直在住院\Nbecause of an immune deficiency.

354
0:19:31.46 --> 0:19:32.61
那你更喜欢哪个\NSo

355
0:19:34.44 --> 0:19:35.76
画画还是做剪贴\NThe painting or the scrapbooking?

356
0:19:37.28 --> 0:19:38.21
做剪贴\NThe scrapbooking.

357
0:19:49.08 --> 0:19:50.31
谁都可以参加吗\NSo anyone can join?

358
0:19:52.89 --> 0:19:55.18
要占用很多课余时间\NIt's a lot of after-school hours

359
0:19:55.18 --> 0:19:58.75
你们有宵禁令  也不准出隔离区\Nand with your curfew and not being permitted outside the Sector...

360
0:19:58.75 --> 0:19:59.57
我是说  不准私自出隔离区\NI mean

361
0:19:59.58 --> 0:20:01.44
不用说了  我懂  我懂\NNo

362
0:20:02.39 --> 0:20:04.63
我的记忆存储在另一台电脑上\NMy memories are stored on a separate computer anyway

363
0:20:04.63 --> 0:20:06.57
我可能不适合做剪贴\Nso I probably wouldn't be good at scrapbooking.

364
0:20:08.69 --> 0:20:09.38
开玩笑的\NJoking.

365
0:20:11.52 --> 0:20:12.36
还是谢了\NThanks anyway.

366
0:20:22.19 --> 0:20:24.09
喜欢搭讪地球女孩啊  歪星人\NYou like talking to our girls there

367
0:20:24.09 --> 0:20:25.42
很漂亮吧\NYou think they're pretty?

368
0:20:26.58 --> 0:20:29.25
我是想加入她的社团  她没同意\NI was interested in joining her club

369
0:20:30.11 --> 0:20:31.28
不用替我遗憾\NNo

370
0:20:31.28 --> 0:20:32.76
还好没成\NIt's just as well.

371
0:20:32.76 --> 0:20:35.02
我也没什么艺术骨头[细胞]\NI don't have an artistic bone in my body.

372
0:20:35.02 --> 0:20:36.82
你全身上下有骨头吗\NDo you have any bones in your body?

373
0:20:36.82 --> 0:20:37.56
说得好\NGood point.

374
0:20:37.56 --> 0:20:38.98
我们的骨骼系统\NOur

375
0:20:38.98 --> 0:20:40.77
由一万条蜈蚣手拉手组成\Nof 10

376
0:20:40.77 --> 0:20:43.51
不过那是奥特里安人的145号秘密\Nbut that's Atrian secret number 145.

377
0:20:43.51 --> 0:20:44.64
你怎么知道的\NHow do you know that?

378
0:20:45.46 --> 0:20:47.29
-一切都好吗  罗曼  -没事\N- Everything okay

379
0:20:48.04 --> 0:20:49.76
我们只是在互相了解\NWe're just getting to know each other.

380
0:20:49.77 --> 0:20:52.98
我不在乎你爸爸是奥特里安人里的大人物\NI don't care that your dad is some big Atrian honcho.

381
0:20:52.98 --> 0:20:54.77
你在这里就要遵守我们的规矩\NYou play by our rules here.

382
0:20:54.77 --> 0:20:55.73
明白了吗  怪胎\NGot it

383
0:20:55.73 --> 0:20:56.95
-得了吧  -泰瑞\N- Come on... - Teri.

384
0:21:01.23 --> 0:21:02.28
明白了\NGot it.

385
0:21:08.03 --> 0:21:09.36
等等  他们因为你跟\NSo

386
0:21:09.37 --> 0:21:10.98
奥特里安人说话而大发雷霆\Nbecause you talked to an Atrian?

387
0:21:12.20 --> 0:21:14.91
我想这有点像是所有人最糟糕的噩梦\NI mean

388
0:21:14.94 --> 0:21:16.18
什么\NWhat is?

389
0:21:16.18 --> 0:21:18.66
人类和奥特里安之间的恋爱倾向\NThe notion of the human-Atrian hookup.

390
0:21:18.66 --> 0:21:20.42
"恋爱"  我们在谈论\NA "Hookup." Okay

391
0:21:20.42 --> 0:21:22.52
绘画和剪贴\Npainting and scrapbooking.

392
0:21:22.53 --> 0:21:24.14
恋爱都是这样开始的\NThat's how it begins.

393
0:21:25.10 --> 0:21:26.68
我认为得给你找一个新的爱好\NI think we need to find you a new hobby

394
0:21:26.68 --> 0:21:29.14
像是吹制玻璃  肚皮舞也可以\Nlike glassblowing

395
0:21:29.14 --> 0:21:30.48
比如...\NSomething like...

396
0:21:37.00 --> 0:21:38.42
这些是什么\NWhat are these?

397
0:21:41.13 --> 0:21:42.45
出院表\NDischarge forms.

398
0:21:43.64 --> 0:21:45.44
你要回家了\NYou're going home?

399
0:21:45.96 --> 0:21:46.98
不回来了\NFor good?

400
0:21:47.72 --> 0:21:48.66
没错\NYeah.

401
0:21:49.45 --> 0:21:50.53
星期二\NTuesday.

402
0:21:51.27 --> 0:21:53.01
你没有跟我说这件事\NA-And you didn't tell me?

403
0:21:58.25 --> 0:21:59.79
化疗不管用\NThe chemo isn't working.

404
0:22:00.62 --> 0:22:01.74
而且找到一个匹配的骨髓\NAnd finding a bone marrow match

405
0:22:01.74 --> 0:22:03.05
也需耗费数年\Ncould take years.

406
0:22:05.57 --> 0:22:07.83
那么你为什么要回家\NSo

407
0:22:07.83 --> 0:22:08.69
我有贫血症\NI'm anemic.

408
0:22:08.69 --> 0:22:10.72
我每隔一周就会得肺炎\NI'm getting pneumonia every other week.

409
0:22:10.74 --> 0:22:12.29
我的血小板数目不停减少\NMy platelet count keeps dropping.

410
0:22:12.30 --> 0:22:13.74
小朱  你不能就这么放弃了\NJules

411
0:22:13.74 --> 0:22:15.25
这不是放弃的问题\NIt's not about giving up.

412
0:22:15.72 --> 0:22:17.29
我一点都不在乎做检查\NI don't even care about the poking

413
0:22:17.29 --> 0:22:19.16
针扎还有呕吐\Nand the prodding and the puking.

414
0:22:20.49 --> 0:22:22.04
我再也承受不了的是\NWhat I can't take anymore

415
0:22:22.05 --> 0:22:24.36
看见我的父母重燃希望\Nis seeing my parents get their hopes up...

416
0:22:24.79 --> 0:22:26.01
却因为一次次的\Nonly to have them dashed

417
0:22:26.01 --> 0:22:28.45
全血细胞数目测试结果而破灭\None lousy CBC test after another.

418
0:22:34.88 --> 0:22:36.66
他们允许你外出一小时的  对吧\NThey'll let you out for an hour

419
0:22:37.44 --> 0:22:38.95
怎么了  我们要去哪里\NWhy? Where are we going?

420
0:22:58.44 --> 0:23:00.31
你确定你知道这个柏草长什么样吗\NAnd you're sure you know what this cyper herb looks like?

421
0:23:00.31 --> 0:23:02.20
-确定  -很好\N- Yes. - Good.

422
0:23:02.20 --> 0:23:05.16
因为我爸绝对不会轻饶我偷了他的通行证\N'Cause my dad will kill me for stealing his access badge.

423
0:23:12.25 --> 0:23:14.63
隔离区执法队\N雷·怀特希尔

424
0:24:06.41 --> 0:24:07.99
[外星语]

425
0:24:08.89 --> 0:24:10.78
[外星语]

426
0:24:11.98 --> 0:24:14.99
要手镯吗\NAh. Bracelets?

427
0:24:14.99 --> 0:24:16.54
-好吧  我买三个  -不用了  谢谢\N- Okay. I'll take three. - Oh

428
0:24:16.54 --> 0:24:18.35
-朱莉娅  不要买  -不用了  谢谢\N- Julia

429
0:24:18.35 --> 0:24:20.42
看看这是谁来了\NWhat do we have here?

430
0:24:22.75 --> 0:24:25.78
你们这些孩子知道不该来到这里的\NYou kids know you're not supposed to be here.

431
0:24:28.03 --> 0:24:29.42
你有柏草吗\NDo you have any cyper?

432
0:24:30.00 --> 0:24:31.21
柏草\NCyper.

433
0:24:34.19 --> 0:24:34.78
跟我来\NRight here.

434
0:24:34.78 --> 0:24:35.83
走吧\NCome on.

435
0:24:36.09 --> 0:24:38.17
[外星语]

436
0:24:38.17 --> 0:24:41.05
[外星语]

437
0:24:48.89 --> 0:24:50.08
一级防范禁闭\NLockdown.

438
0:24:50.35 --> 0:24:53.78
此隔离区现在正在封锁\NThis Sector is now under lockdown!

439
0:24:54.34 --> 0:24:54.93
快走  快走  快走\NCome on

440
0:24:54.93 --> 0:24:56.59
-跟我走  -好的\N- Let's go. - Okay.

441
0:24:56.80 --> 0:24:57.96
一级防范禁闭\NLockdown!

442
0:25:01.13 --> 0:25:02.49
我们在这里是安全的\NWe'll be safe up here.

443
0:25:02.83 --> 0:25:04.50
小朱  你还好吗\NJules

444
0:25:04.70 --> 0:25:05.86
你在跟我开玩笑吗\NAre you kidding?

445
0:25:06.29 --> 0:25:07.88
这真是太不可思议了\NThat was amazing.

446
0:25:08.11 --> 0:25:09.47
你们俩来这里干什么\NWhat are you guys doing here?

447
0:25:11.08 --> 0:25:12.66
我们来找你们的柏草\NWe came for your cyper.

448
0:25:15.32 --> 0:25:17.04
那你为什么不直接说就好了\NWell

449
0:25:17.04 --> 0:25:20.15
我猜它的疗效已经不再只是我们族群之间的秘密了\NI guess its curative powers are no longer just our secret.

450
0:25:22.56 --> 0:25:23.64
快过来\NCome on.

451
0:25:24.98 --> 0:25:27.58
我们从太空船上偷偷带了一些种子进隔离区\NWe smuggled some seeds into the Sector from our ship.

452
0:25:27.58 --> 0:25:29.75
-这是什么地方  -我和我爸一起建造的\N- What is this place? - My dad and I built it.

453
0:25:29.75 --> 0:25:32.35
这儿有点像是我们的私人避难所\NIt's kind of like our own private sanctuary.

454
0:25:32.79 --> 0:25:33.90
就在这里\NHere it is.

455
0:25:38.13 --> 0:25:40.00
天哪  这里有好多柏草\NWow

456
0:25:40.47 --> 0:25:41.93
想拿多少就拿多少吧\NTake as much as you'd like.

457
0:25:42.24 --> 0:25:43.89
这真的是柏草吗\NIs this really cyper?

458
0:25:44.36 --> 0:25:45.36
没错\NYep.

459
0:25:45.89 --> 0:25:47.25
货真价实\NIt's really cyper.

460
0:25:48.67 --> 0:25:50.51
或者是按照你们的说法\NOr

461
0:25:50.51 --> 0:25:51.64
藏红花\Nsaffron.

462
0:25:53.05 --> 0:25:53.79
藏红花\NSaffron?

463
0:25:53.79 --> 0:25:55.06
我们是这么看待它的\NWell

464
0:25:55.06 --> 0:25:56.37
我们用它来做饭\NWe use it for cooking.

465
0:25:56.93 --> 0:25:59.53
但是一些胆大的警卫吹嘘我们香料的神奇功效\NBut some enterprising guards have made a good amount of cash

466
0:25:59.53 --> 0:26:02.08
卖给在隔离区外容易上当受骗的人类\Nselling our spice's magic remedy to...

467
0:26:02.68 --> 0:26:05.09
赚了一大票\Ngullible humans outside the Sector.

468
0:26:07.77 --> 0:26:09.08
这总是值得一试的\NIt was worth a shot.

469
0:26:11.35 --> 0:26:13.69
你有一个患病的朋友吗\NDo you... have a sick friend?

470
0:26:13.69 --> 0:26:15.20
我没有\NUm

471
0:26:16.19 --> 0:26:17.32
她有\NShe does.

472
0:26:23.10 --> 0:26:24.31
对不起\NI'm sorry.

473
0:26:26.09 --> 0:26:28.20
这也不能算是一无所获\NHey

474
0:26:28.29 --> 0:26:30.00
至少我看到了这个地方\NAt least I got to see this place.

475
0:26:30.36 --> 0:26:33.13
我对于奥特里安人的所有东西都很喜欢\NI'm... obsessed with all things Atrian.

476
0:26:34.85 --> 0:26:36.44
你冷得直发抖  给你\NYou're shivering. Here.

477
0:26:39.72 --> 0:26:41.00
谢谢\NThanks.

478
0:26:41.00 --> 0:26:42.45
真有骑士风度\NSo chivalrous.

479
0:26:46.09 --> 0:26:48.00
城市从这上面看好美\NCity looks beautiful from up here.

480
0:26:50.36 --> 0:26:52.33
你这条疤怎么来的\NHow did you get that scar?

481
0:26:53.51 --> 0:26:56.05
这是很久之前的事了\NIt... it happened a long time ago.

482
0:26:57.67 --> 0:26:58.93
我们抵达地球的那天\NOn Arrival Day.

483
0:27:02.17 --> 0:27:05.26
当时在那个棚屋的人是你  对吧\NThat was you i-in the shed

484
0:27:08.25 --> 0:27:10.02
你一直都知道\NYou knew this whole time.

485
0:27:10.35 --> 0:27:11.64
怎么知道的\NHow?

486
0:27:14.25 --> 0:27:15.27
我在学校见到你的那一瞬间\NThe moment I saw you at school

487
0:27:15.27 --> 0:27:16.26
我就产生了这种\NI had this...

488
0:27:16.79 --> 0:27:18.92
对冷掉的意大利面条的异乎寻常的渴望\Nbizarre craving for cold spaghetti.

489
0:27:21.07 --> 0:27:22.65
但是我很确定你已经死了\NBut I was sure you were dead.

490
0:27:22.65 --> 0:27:24.45
我亲眼看到他们射死了你\NI saw them carry you away.

491
0:27:24.45 --> 0:27:25.73
是差点死掉\NI came close.

492
0:27:27.19 --> 0:27:28.43
我的一个心脏确实停跳了\NOne of my hearts actually stopped beating

493
0:27:28.43 --> 0:27:29.78
几分钟\Nfor a few minutes.

494
0:27:30.39 --> 0:27:31.83
幸好我还有备用的\NLuckily

495
0:27:36.68 --> 0:27:40.35
我一直没机会谢谢你当年救了我\NI never got to thank you... for saving my life.

496
0:27:41.54 --> 0:27:43.22
我当时才六岁\NI-I was six.

497
0:27:44.04 --> 0:27:45.23
我没帮到你什么\NI hardly did anything.

498
0:27:45.23 --> 0:27:48.22
别人冷酷无情时  你却宽容相待\NYou were kind... when everyone else was cruel.

499
0:27:48.22 --> 0:27:49.67
这对我意义重大\NThat's something.

500
0:27:53.38 --> 0:27:55.03
宵禁警告\NCurfew warning...

501
0:27:55.03 --> 0:27:57.28
十分钟后就宵禁了  你们该走了\NCurfew starts in ten minutes; you guys should go.

502
0:27:57.57 --> 0:28:00.47
宵禁警告  十...\NCurfew warning: ten...

503
0:28:01.16 --> 0:28:03.41
我不敢相信当年棚屋里的男孩就是他\NI can't believe that was him

504
0:28:03.42 --> 0:28:05.84
不过他显然已经长大了\NBut

505
0:28:06.47 --> 0:28:07.71
朱莉娅\NJulia...

506
0:28:10.50 --> 0:28:11.76
我很抱歉\NI'm sorry.

507
0:28:19.61 --> 0:28:20.96
你在隔离区肯定\NCome on

508
0:28:20.96 --> 0:28:22.67
有男友吧\Nback in the Sector

509
0:28:23.14 --> 0:28:25.28
我真没有\NNo. There's nobody.

510
0:28:25.62 --> 0:28:27.85
谁会想到奥特里安人能这么性感\NWho knew an Atrian could be so sexy?

511
0:28:27.85 --> 0:28:29.29
你是校园里的话题女王\NYou're the talk of the school.

512
0:28:29.29 --> 0:28:30.83
不过别告诉你的朋友泰瑞\NOh

513
0:28:30.83 --> 0:28:32.98
因为我们觉得她可能会咬人\N'cause we think she might bite.

514
0:28:34.78 --> 0:28:36.23
我该去上课了\NI should get to class.

515
0:28:36.82 --> 0:28:38.07
你去哪里  你去哪里\NWhere you going? Where you going?

516
0:28:38.07 --> 0:28:39.33
跟我说说\NHey

517
0:28:39.33 --> 0:28:41.00
有关歪星人的传言是真的吗\Nis it true what they say about Tatties?

518
0:28:41.00 --> 0:28:41.90
别碰我\NOh... Don't!

519
0:28:42.19 --> 0:28:43.59
你性子还挺烈的嘛\NWell

520
0:28:43.59 --> 0:28:45.84
-别碰我  -喂  你\N- Stop it! - Hey! Hey!

521
0:28:46.02 --> 0:28:47.33
离她远点\NLeave her alone.

522
0:28:52.07 --> 0:28:54.04
你非得碍我的事是吧\NYou just can't seem to get out of my way

523
0:28:54.04 --> 0:28:55.41
她是我妹妹\NThat's my sister.

524
0:28:56.99 --> 0:28:57.96
那好\NOkay.

525
0:28:59.63 --> 0:29:00.61
罗曼\NRoman!

526
0:29:01.80 --> 0:29:02.89
罗曼\NRoman!

527
0:29:03.89 --> 0:29:05.04
-来人啊  -放开他\N- Help! - Get off him!

528
0:29:05.04 --> 0:29:05.74
放开他\NHey

529
0:29:05.74 --> 0:29:07.38
谁来帮帮他\NSomebody help him!

530
0:29:09.21 --> 0:29:12.99
别打了  别再打了\NHey! Hey

531
0:29:14.21 --> 0:29:15.29
滚吧\NGet out of here!

532
0:29:15.91 --> 0:29:18.15
滚回你的垃圾星球去吧\NGo back to your own damn planet!

533
0:29:18.38 --> 0:29:20.39
快让他滚\NHey

534
0:29:39.64 --> 0:29:41.87
我观察你很久了  罗曼\NI've been watching you for a long time

535
0:29:42.09 --> 0:29:43.51
我和你父亲一起\NI worked closely with your father

536
0:29:43.51 --> 0:29:45.22
促成了这个项目\Nto help build this program.

537
0:29:45.22 --> 0:29:46.50
你要是不想让它失败\NUnless you want to see it fail

538
0:29:46.51 --> 0:29:48.74
就别像个野兽一样为所欲为\Nyou have to stop acting like an animal.

539
0:29:49.55 --> 0:29:51.86
他攻击我妹妹  我却成了野兽\NHe assaulted my sister

540
0:29:51.86 --> 0:29:53.25
要是这种事再次发生\NThis happens again

541
0:29:53.31 --> 0:29:54.96
他们就会把你交给军医\Nthey'll hand you over to military doctors to have

542
0:29:54.96 --> 0:29:56.65
用你的睾丸做实验\Nyour testicles put in a jar and studied.

543
0:29:56.65 --> 0:29:57.98
你不想这样吧\NYou don't want that.

544
0:29:59.36 --> 0:30:00.38
我也不想这样\NI don't want that.

545
0:30:00.38 --> 0:30:03.03
我的睾丸也绝对不想\NAnd my testicles definitely don't want that.

546
0:30:04.44 --> 0:30:06.84
我一向倾力保护我的投入\NI always protect my investments.

547
0:30:10.00 --> 0:30:12.15
但前提是他们有生存能力\NThat is

548
0:30:23.36 --> 0:30:24.38
罗曼\NRoman.

549
0:30:32.23 --> 0:30:33.10
今晚你会来吧\NYou're coming tonight

550
0:30:33.10 --> 0:30:35.98
没错  格雷森  你会看到她的\NYes

551
0:30:35.98 --> 0:30:37.93
拜托  能洒脱一点吗\NGod. Attempt to be cool.

552
0:30:37.93 --> 0:30:40.41
连像埃默里这样的社交新人\NEven social newbies like Emery

553
0:30:40.41 --> 0:30:42.96
都还是喜欢那些暗恋者不要太过主动\Nstill like a little cool with their drool.

554
0:30:57.14 --> 0:30:59.15
C区清人检查\NQuadrant C cleared for inspection.

555
0:30:59.15 --> 0:30:59.94
不用你管  泰瑞\NStay out of it

556
0:30:59.94 --> 0:31:02.31
就算你能闯进他们的烂派对\NEven if you can get out to crash their lame party

557
0:31:02.31 --> 0:31:04.46
你也不可能在宵禁之前赶回来\Nyou'll never get back in time to make curfew.

558
0:31:04.55 --> 0:31:05.51
不试试看怎么知道\NMaybe

559
0:31:05.51 --> 0:31:07.05
我们得给他们点颜色看看\NBut we need to send a message

560
0:31:07.05 --> 0:31:09.55
不能让今天罗曼的事再重演\NWhat happened to Roman today can never happen again.

561
0:31:13.66 --> 0:31:15.02
都开始穿高跟鞋了\NHigh heels now?

562
0:31:15.17 --> 0:31:16.75
要去约会了\NIt's the beginning of the end.

563
0:31:16.86 --> 0:31:18.70
去参加派对而已\NNo

564
0:31:22.46 --> 0:31:23.66
多加小心\NYou'll be careful?

565
0:31:23.94 --> 0:31:25.02
一向如此\NAlways.

566
0:31:26.18 --> 0:31:27.24
玩得开心点\NHave fun.

567
0:31:27.24 --> 0:31:28.43
你该好好放松一下了\NYou deserve it.

568
0:31:29.63 --> 0:31:31.71
我要去隔离区了\NI'm out of here. Off to the Sectoec

569
0:31:31.71 --> 0:31:32.86
老爸再见\NBye

570
0:31:44.88 --> 0:31:46.12
你想念它吗\NDo you ever miss it?

571
0:31:47.17 --> 0:31:48.15
想念什么\NWhat?

572
0:31:48.81 --> 0:31:49.98
奥特里安吗\NAtria?

573
0:31:50.29 --> 0:31:51.35
当然\NOf course.

574
0:31:51.97 --> 0:31:53.24
没有一天不想\NEvery day.

575
0:31:54.46 --> 0:31:55.46
但是\NBut...

576
0:31:56.24 --> 0:31:58.36
如今这里才是我们的家\Nthis is our home now.

577
0:31:59.01 --> 0:31:59.78
武装守卫\NArmed guards

578
0:31:59.79 --> 0:32:02.40
强制宵禁  还有铁丝网\Nmandatory curfew and barbed wire.

579
0:32:02.61 --> 0:32:04.70
我印象里的家可不是这样的\NIt's not exactly my idea of home.

580
0:32:04.70 --> 0:32:06.06
不会一成不变的\NIt won't be like that forever.

581
0:32:06.07 --> 0:32:07.29
老爸  他们\NDad

582
0:32:07.91 --> 0:32:09.91
他们不会平等对待我们的\Nthey're never gonna treat us like equals.

583
0:32:10.16 --> 0:32:12.28
你们这代人会消除隔阂的\NYou know

584
0:32:12.28 --> 0:32:13.89
有隔阂也是他们造成的\NA gap they created.

585
0:32:15.01 --> 0:32:17.22
我不是说这一切能立马实现\NI-I'm not saying it's gonna happen tomorrow.

586
0:32:17.29 --> 0:32:19.25
这需要时间\NThese things take time.

587
0:32:19.72 --> 0:32:21.57
但我相信他们\NBut I have faith in them.

588
0:32:24.17 --> 0:32:26.58
就如同我相信你一样\NJust like I have faith in you.

589
0:32:27.94 --> 0:32:29.70
你跟其他人不一样\NYou're not like the rest.

590
0:32:31.69 --> 0:32:33.27
我看不到希望\NI just don't see it.

591
0:32:40.43 --> 0:32:42.88
我们明天再接着种吧\NWe'll finish this planting tomorrow.

592
0:32:51.19 --> 0:32:52.40
怎么了  泰瑞\NWhat do you want

593
0:32:52.74 --> 0:32:54.17
我看到你昨晚\NI saw you last night

594
0:32:54.31 --> 0:32:55.95
和那个叫埃默里的女孩在一起\Nwith that girl Emery.

595
0:32:58.24 --> 0:32:59.72
我只是帮她脱离困境\NI was keeping her out of trouble.

596
0:32:59.73 --> 0:33:02.52
我不明白你为什么要费心和他们来往\NI don't know why you even bother making connections with them.

597
0:33:02.65 --> 0:33:05.67
他们很快就会被隔离起来了\NSoon

598
0:33:06.45 --> 0:33:07.46
你这话是什么意思\NWhat are you talking about?

599
0:33:07.46 --> 0:33:08.63
我妈和其他人\NJust something my mother

600
0:33:08.63 --> 0:33:10.81
在他们的周一聚会上这么说的\Nand the rest of the Trags talk about in their Monday meeting.

601
0:33:10.81 --> 0:33:13.27
无意冒犯  泰瑞  但你妈屁也不懂\NNo offense

602
0:33:15.60 --> 0:33:17.06
顺便说一句\NOh

603
0:33:17.68 --> 0:33:19.88
德雷克他们去找那个埃里克了\NDrake and them went after that kid Eric.

604
0:33:24.28 --> 0:33:26.80
一级防范禁闭  2200小时\NLockdown

605
0:33:26.80 --> 0:33:29.51
一级防范禁闭  2200小时\NLockdown

606
0:33:33.59 --> 0:33:35.47
我想要溜出去\NI'm looking for the path to Hades.

607
0:33:37.67 --> 0:33:39.44
你不是诺克斯的儿子吗\NAren't you Nox's son?

608
0:33:45.19 --> 0:33:46.70
你一定是认错人了\NYou must be confused.

609
0:33:47.21 --> 0:33:48.36
能帮我吗\NCan you help?

610
0:33:56.98 --> 0:33:58.16
就像我跟你朋友说的那样\NAs I told your friends

611
0:33:58.16 --> 0:33:59.45
在标记重新启动之前\Nyou've got three hours

612
0:33:59.45 --> 0:34:01.66
你只有三个小时\Nbefore the signal reactivates.

613
0:34:01.66 --> 0:34:04.18
B区8号走廊\NZone B. Corridor 8.

614
0:34:05.34 --> 0:34:07.15
找到那盏坏掉的灯\NLook for the broken light.

615
0:34:07.15 --> 0:34:09.79
正下方有一块假墙板\NDirectly beneath

616
0:34:38.56 --> 0:34:39.46
好了\NOkay.

617
0:34:39.46 --> 0:34:41.09
我不得不动粗才从一个新生那里搞到\NI had to strong-arm a freshman to get the last

618
0:34:41.09 --> 0:34:42.93
最后一点蓝色潘趣\Nof the blue punch.

619
0:34:43.33 --> 0:34:43.96
谢谢\NThanks.

620
0:34:43.96 --> 0:34:44.97
干杯\NCheers.

621
0:34:46.19 --> 0:34:47.30
干杯\NCheers.

622
0:34:48.86 --> 0:34:49.91
你还好吧\NYou all right?

623
0:34:50.65 --> 0:34:51.85
当然了\NYeah. Sure.

624
0:34:51.86 --> 0:34:53.98
你看上去有些心不在焉\NYou just seem a little distracted.

625
0:34:55.94 --> 0:34:57.57
这一切对我来说\NThis is all just...

626
0:34:57.75 --> 0:34:58.86
太新鲜了\Npretty new to me.

627
0:34:58.86 --> 0:35:01.10
要试着放松  没人会说你闲话的\NJust try to relax. You're not being judged.

628
0:35:01.10 --> 0:35:02.54
这里的大多数人都自恋得很\NMost people here are too self-obsessed

629
0:35:02.54 --> 0:35:03.74
才不会关心别人呢\Nto care about anyone else.

630
0:35:06.94 --> 0:35:08.50
大家快来看啊\NHey! You guys got to see this!

631
0:35:14.62 --> 0:35:16.68
看看他们对布洛克干了什么\NLook what someone did to Brock!

632
0:35:19.37 --> 0:35:20.49
谁干的\NWho did this?

633
0:35:21.65 --> 0:35:23.02
谁干的\NWho did this?!

634
0:35:23.91 --> 0:35:25.47
我们只想谈谈\NWe just want to have a little talk.

635
0:35:26.85 --> 0:35:27.87
好吧\NAll right.

636
0:35:28.99 --> 0:35:30.37
今晚得做个了结\NThis ends tonight.

637
0:35:31.20 --> 0:35:33.66
只此一次  我同意\NFor once

638
0:35:37.76 --> 0:35:39.04
我们受够了\NWe're sick of it!

639
0:35:39.50 --> 0:35:41.27
你们这群人每天都来找我们麻烦\NYou guys messing with us every day!

640
0:35:41.27 --> 0:35:44.11
都结束了  如果你们再纠缠我们\NIt's done! It's over! This is what's gonna happen

641
0:35:44.11 --> 0:35:46.20
这就是你们的下场  听明白了吗\Nif you keep messing with us

642
0:35:46.32 --> 0:35:47.97
住手\NHey! Come on!

643
0:35:47.98 --> 0:35:49.37
罗曼  你来干什么\NRoman

644
0:35:49.37 --> 0:35:50.11
你知道他们发现你们出了隔离区\NDo you know what they'd do to you

645
0:35:50.11 --> 0:35:51.26
会拿你怎么样吗\Nif they found you outside the Sector?

646
0:35:51.26 --> 0:35:52.63
你应该和我一条战线\NYou should be fighting alongside of me

647
0:35:52.63 --> 0:35:54.10
而不是和我对着干\Ninstead of against me!

648
0:35:54.30 --> 0:35:55.15
警察\NCops!

649
0:35:55.16 --> 0:35:56.25
快跑\NCome on!

650
0:35:59.40 --> 0:36:01.79
罗曼  罗曼  罗...\NRoman. Roman! Ro...!

651
0:36:05.81 --> 0:36:08.11
-你没事吧  -埃默里  快走\N- Are you okay? - Emery

652
0:36:08.12 --> 0:36:10.43
快来  罗曼  上车\NCome on! Roman

653
0:36:10.43 --> 0:36:12.68
你疯了吗  埃默里  别管他了\NAre you nuts? Come on

654
0:36:12.68 --> 0:36:13.53
警察来了怎么办\NBut what about the cops?!

655
0:36:13.53 --> 0:36:15.37
那是他的问题  不关我们的事  快走\NHis problem

656
0:36:15.37 --> 0:36:17.19
埃默里  你要干什么  走啊\NEmery

657
0:36:17.90 --> 0:36:19.63
-埃默里  -格雷森  我们必须离开\N- Emery! - Grayson

658
0:36:19.63 --> 0:36:21.21
罗曼  快跑  跑比较快\NRoman

659
0:36:21.21 --> 0:36:22.77
-我们不能丢下她  -当然可以\N- We can't just leave her! - Yes

660
0:36:22.77 --> 0:36:23.80
她简直是个疯子\NShe's a total freak!

661
0:36:23.80 --> 0:36:25.56
伙计  她选择了那个歪星人\NDude

662
0:36:25.56 --> 0:36:26.66
认了吧\NOwn it.

663
0:36:26.67 --> 0:36:28.34
快点  我们走\NCome on! Let's go!

664
0:36:51.84 --> 0:36:52.94
你还好吗\NAre you okay?

665
0:36:53.02 --> 0:36:56.10
除了刚才社交自杀了一把\NOther than having possibly just committed social suicide?

666
0:36:53.02 --> 0:36:56.10
社交自杀  指躲开人群  避免社交活动

667
0:36:56.96 --> 0:36:57.98
还好\NSure.

668
0:37:00.85 --> 0:37:01.90
你呢\NYou okay?

669
0:37:05.39 --> 0:37:06.66
真好笑\NYou know

670
0:37:07.12 --> 0:37:08.48
这么多年来我父亲一直在教我\NAll these years

671
0:37:08.48 --> 0:37:11.16
那被称道的所谓"人性"\Nabout this vaunted thing called "Humanity

672
0:37:12.18 --> 0:37:13.36
从定义上来说  那是我们\Nsomething that

673
0:37:13.36 --> 0:37:15.50
不可能拥有的\Nwe could never possess.

674
0:37:16.46 --> 0:37:17.51
但是...\NBut...

675
0:37:18.33 --> 0:37:21.64
和你们人类接触了几天之后\Nafter spending just a few days amongst your kind...

676
0:37:23.02 --> 0:37:25.84
我弄不明白"人性"这个概念了\Nthis concept of "Humanity" doesn't seem so clear to me.

677
0:37:25.84 --> 0:37:29.28
不是的  世上还是有很多好人的\NNo. There are a lot of good people out there.

678
0:37:29.28 --> 0:37:30.20
是吗\NYeah?

679
0:37:31.26 --> 0:37:33.18
我只遇见过一个\NWell

680
0:37:56.13 --> 0:37:57.80
你不必害怕\NYou don't have to be afraid.

681
0:38:08.47 --> 0:38:09.65
对不起\NSorry.

682
0:38:10.10 --> 0:38:11.55
是我妈\NIt's my mom.

683
0:38:13.24 --> 0:38:14.55
我一会儿就到家\NI'll be home soon.

684
0:38:18.94 --> 0:38:21.27
好的  我正在路上\NO-Okay. I'm on my way.

685
0:38:23.29 --> 0:38:25.39
是...朱莉娅\NIt's... Julia.

686
0:38:26.27 --> 0:38:27.79
我得走了\NI have to go.

687
0:38:29.28 --> 0:38:30.06
我明白\NI understand.

688
0:38:30.06 --> 0:38:31.12
继续往前直走\NIf you keep walking straight

689
0:38:31.12 --> 0:38:32.98
你就会找到回隔离区的路\Nyou'll find the road to the Sector.

690
0:38:33.16 --> 0:38:34.19
谢谢你\NThanks.

691
0:38:38.44 --> 0:38:39.55
你没事吧\NAre you okay?

692
0:38:40.29 --> 0:38:41.73
我得走了\NI have to go.

693
0:39:05.94 --> 0:39:06.86
泰瑞\NTeri.

694
0:39:07.00 --> 0:39:08.10
几小时前你和罗曼在一起\NYou were with Roman a few hours ago.

695
0:39:08.10 --> 0:39:10.10
你知道他现在在哪里吗  都快宵禁了\NDo you know where he is? It's almost curfew.

696
0:39:11.03 --> 0:39:12.05
泰瑞\NTeri?

697
0:39:56.79 --> 0:39:58.91
你怎么能让我儿子出去\NHow could you let my son go out there?

698
0:40:09.30 --> 0:40:10.41
先生们  出什么事了\NIs there a problem here

699
0:40:10.41 --> 0:40:11.88
没事\NNo

700
0:40:17.57 --> 0:40:19.69
离开桌子  让我看到你们的手\NStep away from the table. Let me see your hands.

701
0:40:26.79 --> 0:40:27.94
那是什么\NWhat's that?

702
0:40:27.97 --> 0:40:29.33
这可不行  威布尔\NNot good

703
0:40:37.39 --> 0:40:38.75
放下武器\NDrop your weapon!

704
0:40:43.01 --> 0:40:44.17
威布尔  不要\NWeeble

705
0:41:01.99 --> 0:41:03.09
有人中枪\NMan down!

706
0:41:09.58 --> 0:41:11.02
怀特希尔指挥官\NCommander Whitehill!

707
0:41:12.37 --> 0:41:13.55
雷\NRay!

708
0:41:14.77 --> 0:41:16.03
你没事吧\NAre you okay?

709
0:41:50.00 --> 0:41:56.21
星恋`)
	if !(len(lang) == 2 && lang[0]=="zh" && lang[1]=="en") {
		t.Errorf("Expect zh,en but %v", lang)
	}
}




func TestCld(t *testing.T) {
	lang := DetectLanguage(`1
00:00:02,050 --> 00:00:35,550
【造雨人】

2
00:00:37,550 --> 00:00:41,158
我父亲一生都恨律师

3
00:00:41,159 --> 00:00:42,958
我个人认为,他不算是个好人.

4
00:00:42,959 --> 00:00:45,388
他汹酒并且虐待我的母亲.

5
00:00:45,388 --> 00:00:47,257
他也殴打我.

6
00:00:47,257 --> 00:00:49,226
要是你认为我要成为一个律师

7
00:00:49,226 --> 00:00:50,965
只是为了气他,

8
00:00:50,966 --> 00:00:52,995
那你也许就错了.

9
00:00:52,995 --> 00:00:55,094
我之所以想成为一个律师

10
00:00:55,094 --> 00:00:56,803
是自从我读了有关民权律师

11
00:00:56,804 --> 00:00:58,503
在五十和六十年代

12
00:00:58,503 --> 00:01:00,902
为法律发掘出令人惊异的运用.

13
00:01:00,902 --> 00:01:03,141
他们做了许多人认为

14
00:01:03,141 --> 00:01:04,840
本来是不可能的事.

15
00:01:04,841 --> 00:01:07,410
他们给了律师一个好名声.

16
00:01:07,410 --> 00:01:10,109
因此我进了法学院,

17
00:01:10,109 --> 00:01:11,978
这确实激怒了我的父亲,

18
00:01:11,978 --> 00:01:13,907
但无论如何他是被气过了.

19
00:01:13,908 --> 00:01:15,977
在我读大学的第一年,有一天他喝醉了

20
00:01:15,977 --> 00:01:17,046
从梯子上摔下来,

21
00:01:17,047 --> 00:01:19,386
这发生在他干活的工厂里,

22
00:01:19,386 --> 00:01:21,915
而他却不知道要先告谁.

23
00:01:21,915 --> 00:01:24,414
大约两个月后他就因此去世了.

24
00:01:24,414 --> 00:01:26,923
"我喝干啦".

25
00:01:26,923 --> 00:01:28,522
我的一些同学

26
00:01:28,522 --> 00:01:30,691
知道他们将从学校出去后直接

27
00:01:30,692 --> 00:01:32,361
进入顶级的律师事务所,

28
00:01:32,361 --> 00:01:34,990
主要依靠他们的家庭关系网.

29
00:01:34,990 --> 00:01:36,259
而我所仅有的一个关系

30
00:01:36,260 --> 00:01:38,859
是三年来在一个酒巴跑堂而建立起来的,

31
00:01:38,859 --> 00:01:40,098
因为我得赚钱付我的学费.

32
00:01:40,098 --> 00:01:41,997
我依然有计划要使司法

33
00:01:41,998 --> 00:01:43,627
发出耀眼光芒

34
00:01:43,627 --> 00:01:44,696
照亮每一个黑暗角落,

35
00:01:44,697 --> 00:01:47,635
但我目前真正需要的只是一个工作,

36
00:01:47,636 --> 00:01:48,765
非常迫切,

37
00:01:48,765 --> 00:01:52,573
因为在孟菲斯已经有太多的律师了.

38
00:01:55,173 --> 00:01:57,372
这座城市大受其扰.

39
00:01:57,372 --> 00:02:01,870
《造 雨 人》

40
00:02:01,871 --> 00:02:03,640
《造 雨 人》
"我不认为是这样."

41
00:02:03,640 --> 00:02:06,648
《造 雨 人》
我想像不出还有什么事能比
这更令人苦恼了,

42
00:02:06,649 --> 00:02:07,908
那就是告诉别人我所为之工作的人

43
00:02:07,909 --> 00:02:12,917
是伯鲁瑟.斯通.

44
00:02:12,917 --> 00:02:17,015
我是说,他是个律师,大家叫他伯鲁瑟.

45
00:02:17,015 --> 00:02:19,684
这就是为什么我如此失望.

46
00:02:19,684 --> 00:02:20,723
"当然."

47
00:02:20,724 --> 00:02:22,693
把门关上.

48
00:02:22,693 --> 00:02:25,661
"不,不是那个,也不是那个."

49
00:02:25,662 --> 00:02:29,430
"当然."

50
00:02:29,431 --> 00:02:32,299
"到我家里来做这件事吧."

51
00:02:32,300 --> 00:02:35,998
"好的."

52
00:02:35,999 --> 00:02:38,568
我真的同情这个可怜的
联邦调查局调查员

53
00:02:38,568 --> 00:02:40,367
他非要想从交谈中

54
00:02:40,367 --> 00:02:41,836
搞些资料出来.

55
00:02:41,837 --> 00:02:44,136
律师的办公室里养着活生生的鲨鱼.

56
00:02:44,136 --> 00:02:46,135
这是个玩笑.
明白吗?

57
00:02:46,135 --> 00:02:47,634
如此说来,普林斯,

58
00:02:47,634 --> 00:02:49,543
这就是你向我提起过的法学院学生吗?

59
00:02:49,544 --> 00:02:50,773
是,先生.我刚在

60
00:02:50,773 --> 00:02:52,612
孟菲斯州立大学读完第三年.

61
00:02:52,613 --> 00:02:54,142
你能在这里雇他吗?

62
00:02:54,142 --> 00:02:56,141
你看,我可以为他担保.

63
00:02:56,141 --> 00:02:57,640
这孩子需要找份工作.

64
00:02:57,641 --> 00:03:01,349
他在我们俱乐部的酒巴打工.

65
00:03:01,350 --> 00:03:04,518
鲁迪,这儿会是你工作的好地方.

66
00:03:04,519 --> 00:03:06,248
这儿会是...

67
00:03:06,248 --> 00:03:08,147
你工作的好地方.

68
00:03:08,147 --> 00:03:10,716
你可以看到真正的律师是怎么做的.

69
00:03:10,716 --> 00:03:13,984
现在,这还不能算正式的位置.

70
00:03:13,985 --> 00:03:15,124
不算吗?

71
00:03:15,125 --> 00:03:18,923
我希望我的人自己为自己发工资--

72
00:03:18,923 --> 00:03:22,091
赚取他们自己的酬金.

73
00:03:22,092 --> 00:03:23,531
说点儿什么让我听听.

74
00:03:23,532 --> 00:03:25,861
那是如何实施的呢?

75
00:03:25,861 --> 00:03:28,969
嗯,鲁迪,你每月支取一千美金,

76
00:03:28,970 --> 00:03:31,998
另外你可以从你赚取的797F4CD3里提取三分之一.

77
00:03:31,999 --> 00:03:35,167
要是你不能赚的比你支取的多

78
00:03:35,168 --> 00:03:37,337
那么到每个月底你就欠了我那部分差额.

79
00:03:37,337 --> 00:03:39,376
你听明白了吗?

80
00:03:39,376 --> 00:03:41,135
我觉得挺公平,鲁迪.

81
00:03:41,136 --> 00:03:42,905
这是个不易接受的交易,的确如此.

82
00:03:42,905 --> 00:03:45,044
但这么做你能挣很多钱.

83
00:03:45,044 --> 00:03:46,113
嗯...

84
00:03:46,114 --> 00:03:47,743
让我来告诉你,鲁迪.

85
00:03:47,743 --> 00:03:49,542
这是我唯一的经营方式.

86
00:03:49,543 --> 00:03:51,242
我会让你处理很多有利的案子.

87
00:03:51,242 --> 00:03:54,080
我手里有案子.

88
00:03:54,081 --> 00:03:55,510
哦...

89
00:03:55,510 --> 00:03:58,179
我现在手上有两个案子.

90
00:03:58,180 --> 00:03:59,889
一个是办一个遗嘱,

91
00:03:59,889 --> 00:04:01,288
我替一个老妇人起个遗嘱草稿.

92
00:04:01,288 --> 00:04:02,957
她很富有,值好几百万.

93
00:04:02,958 --> 00:04:05,087
而手上的另一个案子,

94
00:04:05,087 --> 00:04:06,326
是个关于保险的案子,

95
00:04:06,327 --> 00:04:09,495
巨大福利保险公司,
我想你听说过这家公司.

96
00:04:09,496 --> 00:04:11,695
你有这些委托人的签字吗?

97
00:04:11,695 --> 00:04:13,894
我现在正要去看他们.

98
00:04:13,894 --> 00:04:15,623
他们也很愿意听我的意见.

99
00:04:15,623 --> 00:04:17,292
我在律师实习处帮过他们的忙.

100
00:04:17,293 --> 00:04:19,002
很好,我会让你去和我的一个同事

101
00:04:19,002 --> 00:04:20,901
狄克.谢佛来尔谈谈.

102
00:04:20,902 --> 00:04:24,530
他通常与那些大的保险公司打交道.

103
00:04:24,530 --> 00:04:27,339
他处理这儿所有具高度权限的事务.

104
00:04:27,339 --> 00:04:29,808
嘿,狄克!

105
00:04:29,808 --> 00:04:31,007
狄克!

106
00:04:31,008 --> 00:04:32,037
该死的.

107
00:04:32,038 --> 00:04:33,577
你建立了那些案子的档案了吗?

108
00:04:33,577 --> 00:04:34,536
在我的车上.

109
00:04:34,537 --> 00:04:35,636
很好.

110
00:04:39,015 --> 00:04:41,544
嘿,帅哥.

111
00:04:41,544 --> 00:04:44,243
这是狄克.谢佛来尔.

112
00:04:44,243 --> 00:04:46,812
这家伙会带你上路.

113
00:04:46,812 --> 00:04:48,411
我想让你做的是

114
00:04:48,412 --> 00:04:51,181
立案起诉这个巨大福利保险公司,

115
00:04:51,181 --> 00:04:53,990
你把我的名字放在上面.

116
00:04:53,990 --> 00:04:55,759
我们今天就建立档案.

117
00:04:55,759 --> 00:04:57,158
谢谢你.

118
00:04:57,159 --> 00:04:59,288
鲁迪,你会学到很多东西.

119
00:04:59,288 --> 00:05:01,887
鲁迪,我很高兴你到这儿来.

120
00:05:01,887 --> 00:05:03,896
你做了个很好的选择,孩子.

121
00:05:03,896 --> 00:05:05,125
不错,那么,谢谢你.

122
00:05:05,126 --> 00:05:06,695
好的,谢谢你.

123
00:05:06,695 --> 00:05:08,394
出去的时候把门带上,好吗?

124
00:05:08,395 --> 00:05:09,894
律师事务所,我能帮你做什么吗?

125
00:05:09,894 --> 00:05:11,333
这是个办公室,有人在里面.

126
00:05:11,334 --> 00:05:12,703
要是有人在办公室里,

127
00:05:12,703 --> 00:05:14,102
你就别用它,它满员了.

128
00:05:14,103 --> 00:05:15,562
这里是厕所.

129
00:05:15,562 --> 00:05:16,571
"等一会儿".

130
00:05:16,572 --> 00:05:18,201
对不起.

131
00:05:18,201 --> 00:05:19,940
那么,你是他的的同事吗?

132
00:05:19,941 --> 00:05:21,540
算是吧.不是正式的.

133
00:05:21,540 --> 00:05:23,009
确凿说来我还不算个律师.

134
00:05:23,010 --> 00:05:25,339
伯鲁瑟一般要求我

135
00:05:25,339 --> 00:05:27,178
审查那些牵涉保险的案子.

136
00:05:27,178 --> 00:05:29,607
我以前一直为特殊的共同保险公司工作.

137
00:05:29,607 --> 00:05:31,406
我厌倦了就进了法学院.

138
00:05:31,407 --> 00:05:33,816
那么,你什么时候从法学院毕业的?

139
00:05:33,816 --> 00:05:35,075
五年前.

140
00:05:35,075 --> 00:05:36,344
瞧,我在律师执照考试上

141
00:05:36,345 --> 00:05:37,584
有那么点儿小麻烦.

142
00:05:37,584 --> 00:05:39,683
我考了六次.

143
00:05:39,684 --> 00:05:40,683
听到这些真让我感到遗憾.

144
00:05:40,683 --> 00:05:41,712
是啊.你什么时候参加执照考试?

145
00:05:41,713 --> 00:05:42,722
我差不多三周后就要考了.

146
00:05:42,723 --> 00:05:43,722
噢,是啊.

147
00:05:43,722 --> 00:05:44,751
这真那么难吗?

148
00:05:44,752 --> 00:05:46,151
我要说,相当地难.

149
00:05:46,151 --> 00:05:47,420
我一年前考过,

150
00:05:47,421 --> 00:05:48,990
我想我是再也不要去考了.

151
00:05:48,990 --> 00:05:50,919
随它去吧,这里是伯鲁瑟的法律图书室.

152
00:05:50,920 --> 00:05:52,159
要是你从冰箱里取任何东西

153
00:05:52,159 --> 00:05:54,158
或者你要用这个冰箱,

154
00:05:54,159 --> 00:05:56,098
你可以在你的东西上写上你的名字.

155
00:05:56,098 --> 00:05:58,797
不过他们照样会吃了它.
替我放一放.

156
00:05:58,797 --> 00:06:00,526
噢,这只是米饭罢了.
这是什么玩意儿.黛茜.

157
00:06:00,526 --> 00:06:01,525
哎.

158
00:06:01,526 --> 00:06:03,335
我们把这儿搞的很糟.清理一下,啊?

159
00:06:03,335 --> 00:06:04,334
当然,狄克.

160
00:06:04,335 --> 00:06:05,564
随它的便...

161
00:06:05,564 --> 00:06:07,763
等等,你上法庭的时候怎么办?

162
00:06:07,764 --> 00:06:10,003
老实说我是很少单独出庭的.

163
00:06:10,003 --> 00:06:11,032
所以至今还没被抓住过.

164
00:06:11,033 --> 00:06:12,202
这里有这么多律师,

165
00:06:12,202 --> 00:06:16,240
要注意上我们是不可能的.

166
00:06:16,241 --> 00:06:18,170
伯鲁瑟拥有这一切.

167
00:06:18,170 --> 00:06:19,909
噢,真不错.

168
00:06:19,909 --> 00:06:22,238
你真不能称其为律师事务所.

169
00:06:22,239 --> 00:06:24,308
每个人都只是为了自己.

170
00:06:24,308 --> 00:06:26,607
你会学到的.

171
00:06:28,147 --> 00:06:30,116
什么,你是在搬家吗?

172
00:06:30,116 --> 00:06:32,315
被赶出来了.

173
00:06:33,555 --> 00:06:35,324
这是保险单.

174
00:06:37,323 --> 00:06:38,352
不错.

175
00:06:38,353 --> 00:06:39,422
你觉的怎么样?

176
00:06:39,423 --> 00:06:41,852
嗯,这就是这行业里

177
00:06:41,852 --> 00:06:42,991
有名无实的玩意儿.

178
00:06:42,991 --> 00:06:44,820
他们拒付的理由是什么?

179
00:06:44,821 --> 00:06:45,890
噢,这个,所有的方面.

180
00:06:45,890 --> 00:06:48,828
他们先在主要方面拒付了,

181
00:06:48,829 --> 00:06:51,897
然后他们说白血病

182
00:06:51,898 --> 00:06:53,027
是在投保之前就有的,

183
00:06:53,028 --> 00:06:54,027
接着又说

184
00:06:54,027 --> 00:06:55,866
白血病根本就不在保险计划之内.

185
00:06:55,867 --> 00:06:57,236
我这里有七封信.

186
00:06:57,236 --> 00:06:59,305
他们付没付所有保险项目的保险金?

187
00:06:59,305 --> 00:07:00,364
根据勃拉克女士的说法,

188
00:07:00,365 --> 00:07:01,834
任何一项保险费她都付过了.

189
00:07:01,835 --> 00:07:03,434
这些个杂种.

190
00:07:03,434 --> 00:07:07,272
这是份标准的黑人们称之为

191
00:07:07,273 --> 00:07:09,742
街头推销保险合同.

192
00:07:09,742 --> 00:07:12,580
我该怎么办?

193
00:07:12,581 --> 00:07:15,979
你签字接办.所有的都签.

194
00:07:15,980 --> 00:07:18,249
然后交给简.李曼.斯通事务所.

195
00:07:18,249 --> 00:07:19,878
好吧.

196
00:07:19,878 --> 00:07:21,947
这就对了.我会在这个案子上帮助你.

197
00:07:21,947 --> 00:07:23,186
好的.那么,谢谢你.

198
00:07:23,187 --> 00:07:24,286
十分感激.

199
00:07:24,287 --> 00:07:26,756
没什么比把保险公司钉在架子上

200
00:07:26,756 --> 00:07:28,955
更令人激动的事了.

201
00:07:46,369 --> 00:07:49,168
嗨,勃拉克女士.我是鲁迪.拜勒.

202
00:07:49,168 --> 00:07:51,767
还记得吗?我在处理你控告

203
00:07:51,767 --> 00:07:52,966
巨大福利保险公司的案子.

204
00:07:52,967 --> 00:07:54,006
我在孟菲斯州立大学

205
00:07:54,006 --> 00:07:56,675
法律实习室和你碰过头.

206
00:07:56,675 --> 00:07:57,974
对,快请进来.

207
00:07:57,975 --> 00:08:01,473
请进.很抱歉那几条疯了的狗.

208
00:08:01,474 --> 00:08:02,843
噢,那没什么.

209
00:08:02,843 --> 00:08:04,942
我还以为你是
"耶和华见证人"的传教士呢.

210
00:08:04,942 --> 00:08:05,941
伯迪在哪儿?

211
00:08:05,942 --> 00:08:07,681
他在外面车里.

212
00:08:07,681 --> 00:08:10,050
他在那里做什么?

213
00:08:10,050 --> 00:08:12,149
他哪儿都不去.

214
00:08:12,150 --> 00:08:13,949
他的脑袋不太对劲儿,是战争创伤,

215
00:08:13,949 --> 00:08:16,917
在朝鲜留下的.你听说过飞机场的

216
00:08:16,918 --> 00:08:18,147
金属探测器吧?

217
00:08:18,148 --> 00:08:20,157
他就是脱光了通过检测门

218
00:08:20,157 --> 00:08:21,456
警铃都会响个不停.

219
00:08:21,456 --> 00:08:23,755
他的脑袋里还留着块弹片.

220
00:08:23,756 --> 00:08:25,555
哦.

221
00:08:25,555 --> 00:08:29,823
这真是可怕.

222
00:08:29,823 --> 00:08:31,892
丹尼.雷还好吗?

223
00:08:31,893 --> 00:08:33,892
嗯,时好时坏的.

224
00:08:33,892 --> 00:08:34,861
"脑袋里留着块弹片".

225
00:08:34,862 --> 00:08:36,661
你要见他吗?

226
00:08:36,661 --> 00:08:38,130
也许等会儿吧.

227
00:08:38,131 --> 00:08:42,169
嗯,现在...

228
00:08:42,169 --> 00:08:47,407
巨大福利保险公司去年八月
当丹尼.雷被诊断之后

229
00:08:47,407 --> 00:08:49,876
第一次拒绝了你的索赔.

230
00:08:49,876 --> 00:08:52,844
为什么一直等到现在你们才向律师咨询?

231
00:08:52,845 --> 00:08:55,004
我猜是由于愚蠢吧.

232
00:08:55,005 --> 00:08:56,514
我只是不断地给他们写信,

233
00:08:56,514 --> 00:08:58,473
而他们也不断地给我回信,

234
00:08:58,473 --> 00:09:03,311
哪,这是最后收到的一封.

235
00:09:03,312 --> 00:09:04,811
"亲爱的勃拉克女士,

236
00:09:04,811 --> 00:09:06,120
以前的七次情况是

237
00:09:06,121 --> 00:09:08,150
我们公司已经书面拒绝了你的索赔要求.

238
00:09:08,150 --> 00:09:11,048
现在我们第八次也是最后一次
拒绝那个要求.

239
00:09:11,049 --> 00:09:16,027
你一定是很愚蠢,很愚蠢,相当的愚蠢."

240
00:09:16,027 --> 00:09:24,364
"真诚地,艾沃特.鲁富肯,
索赔部付总裁."

241
00:09:24,364 --> 00:09:26,493
真令人难以置信.

242
00:09:26,493 --> 00:09:28,732
你就是那位律师.

243
00:09:28,733 --> 00:09:35,600
我妈夸过你好多次.

244
00:09:35,600 --> 00:09:37,639
她说你会控告那些个

245
00:09:37,639 --> 00:09:39,808
巨大福利保险公司的杂种.

246
00:09:39,809 --> 00:09:41,868
让他们掏腰包,是不是?

247
00:09:41,868 --> 00:09:42,977
对,是这样.

248
00:09:42,978 --> 00:09:44,907
是这样的.

249
00:09:44,907 --> 00:09:46,606
嗨,妈妈.

250
00:09:46,606 --> 00:09:47,845
嗨,亲爱的.

251
00:09:47,846 --> 00:09:50,914
好吧,在我们索赔能建档之前,

252
00:09:50,915 --> 00:09:53,184
我需要你们三个人的签字.

253
00:09:53,184 --> 00:09:55,413
爸爸参与进来了吗?

254
00:09:55,413 --> 00:09:57,882
唉,他说他没有.

255
00:09:57,882 --> 00:09:59,451
也许有那么一天他会参与进来,

256
00:09:59,452 --> 00:10:00,621
但另一天他又说不参与.

257
00:10:00,621 --> 00:10:01,850
好吧,这是份合同书.

258
00:10:01,851 --> 00:10:03,690
上面写了些什么?

259
00:10:03,690 --> 00:10:05,419
噢,很平常的.

260
00:10:05,420 --> 00:10:07,089
上面说地很清楚.

261
00:10:07,089 --> 00:10:08,258
基本上讲的是

262
00:10:08,259 --> 00:10:10,888
你们将请我们做为你们的代理人,

263
00:10:10,888 --> 00:10:12,757
我们将接手处理你们的案子.

264
00:10:12,757 --> 00:10:15,396
并负担所有的可能开销,

265
00:10:15,396 --> 00:10:17,725
然后我们从赔偿款中抽取三分之一佣金.

266
00:10:17,725 --> 00:10:19,294
嗯嗯.

267
00:10:19,295 --> 00:10:23,233
好吧,可为什么这几句话要写两页纸?

268
00:10:23,233 --> 00:10:25,862
别点.

269
00:10:25,862 --> 00:10:29,870
怪不得我要死了.

270
00:10:29,871 --> 00:10:32,709
我们三个人都必须签吗?

271
00:10:32,710 --> 00:10:34,839
是的,夫人,三个人都得签.
签在有你名字的地方.

272
00:10:34,839 --> 00:10:36,908
他说他并没有参与进来.

273
00:10:36,909 --> 00:10:38,538
拿枝笔去找他

274
00:10:38,538 --> 00:10:40,207
让他在这个该死的东西上签个字就是了.

275
00:10:40,207 --> 00:10:43,205
我想我只能这么办了.

276
00:10:49,144 --> 00:10:50,253
喂,伯迪,

277
00:10:50,254 --> 00:10:52,183
你必须在文件上签个字

278
00:10:52,183 --> 00:10:54,882
以便丹尼.雷可以得到他所需的治疗,

279
00:10:54,882 --> 00:10:57,291
我可不想因为你而惹出麻烦来.

280
00:10:57,291 --> 00:10:59,220
给我这可恶的酒瓶子

281
00:10:59,221 --> 00:11:01,220
我要把它扔到大街上去.

282
00:11:01,220 --> 00:11:05,328
现在,在这该死的玩意儿上签字.
快点,你手脚快点儿嘛.

283
00:11:06,098 --> 00:11:10,566
我知道你一定以为他们都疯了.

284
00:11:10,567 --> 00:11:17,204
他们都是很好的人.

285
00:11:17,204 --> 00:11:21,332
嗨,哥们儿.

286
00:11:21,333 --> 00:11:25,771
嗨,你留鼻血了.

287
00:11:25,771 --> 00:11:27,610
勃拉克女士!

288
00:11:27,611 --> 00:11:28,940
把头仰起来.

289
00:11:28,940 --> 00:11:31,109
勃拉克女士!

290
00:11:31,109 --> 00:11:32,848
没什么.我来弄.

291
00:11:32,849 --> 00:11:34,148
他在流鼻血.

292
00:11:34,148 --> 00:11:36,117
把头仰起来.仰起头来.

293
00:11:36,118 --> 00:11:37,777
快点儿,仰起头来.

294
00:11:37,777 --> 00:11:41,645
啊,亲爱的,没什么,一切都会好的.

295
00:11:41,646 --> 00:11:43,085
我自己来.

296
00:11:43,085 --> 00:11:45,754
让我自己来.

297
00:11:45,754 --> 00:11:47,283
都会好的,都会好的.

298
00:11:47,284 --> 00:11:48,883
好啦,我捂住了.

299
00:11:48,883 --> 00:11:51,292
你一定会好起来的.

300
00:11:51,292 --> 00:11:52,621
文件在哪里?

301
00:11:52,622 --> 00:11:54,321
丹尼.雷,你不能稍稍歇一歇吗?

302
00:11:54,321 --> 00:11:55,390
你可以等一会儿的.

303
00:11:55,391 --> 00:11:56,360
你用不着急着去做.

304
00:11:56,360 --> 00:11:58,029
不,不.我希望现在就做.

305
00:11:58,030 --> 00:12:00,229
好吧.

306
00:12:00,229 --> 00:12:03,067
对,对,你能做到.

307
00:12:03,068 --> 00:12:04,397
请继续.

308
00:12:04,398 --> 00:12:08,076
就好了.

309
00:12:17,243 --> 00:12:21,371
啊,柏蒂女士,这是鲁迪.拜勒.

310
00:12:21,372 --> 00:12:22,881
什么?

311
00:12:22,881 --> 00:12:25,739
你是谁?

312
00:12:25,740 --> 00:12:27,679
这是鲁迪.拜勒.

313
00:12:27,679 --> 00:12:28,748
噢.

314
00:12:28,749 --> 00:12:30,218
我们在孟菲斯州立大学的

315
00:12:30,218 --> 00:12:31,247
法律实习室见过面.

316
00:12:31,248 --> 00:12:34,186
噢,对了.
噢,请进,快请进.

317
00:12:34,187 --> 00:12:36,256
多谢,多谢.

318
00:12:36,256 --> 00:12:37,455
您今天还好吧?

319
00:12:37,456 --> 00:12:38,615
不错,很不错.

320
00:12:38,615 --> 00:12:40,754
柏蒂女士,我想跟您谈谈

321
00:12:40,755 --> 00:12:41,884
有关您的遗嘱.

322
00:12:41,884 --> 00:12:44,093
把我的孩子们从我的遗嘱里剔出去.
剔除,剔除,剔除.

323
00:12:44,094 --> 00:12:45,493
剔除,剔除,剔除.
我知道了.

324
00:12:45,493 --> 00:12:46,762
我昨晚上没睡好,

325
00:12:46,763 --> 00:12:48,792
原因是我一直在担心您的财产.

326
00:12:48,792 --> 00:12:51,690
现在,要是您不当心的话,柏蒂女士,

327
00:12:51,691 --> 00:12:54,360
政府就会从中抽取相当大的份额.

328
00:12:54,360 --> 00:12:56,969
好吧,相当大一部分税额
假如做一个稍微小心点儿的

329
00:12:56,969 --> 00:12:59,738
遗产计划是可以避开的.

330
00:12:59,738 --> 00:13:01,967
噢,那是冠冕堂皇合法的.

331
00:13:01,967 --> 00:13:03,936
我假设你是希望你的名字

332
00:13:03,937 --> 00:13:05,506
能写进我的遗嘱里去.

333
00:13:05,506 --> 00:13:06,775
当然不是,夫人.

334
00:13:06,776 --> 00:13:08,675
律师们总是想把他们的名字

335
00:13:08,675 --> 00:13:09,844
写进我的遗嘱.

336
00:13:09,844 --> 00:13:12,742
不是这样的,夫人.嗯,律师是多种多样的.

337
00:13:12,743 --> 00:13:15,442
我所需要从您那儿得到的只是
那个遗产计划的意图,

338
00:13:15,442 --> 00:13:17,581
我需要知道钱在什么地方放着.

339
00:13:17,582 --> 00:13:21,310
那是公债,股票还是现金?

340
00:13:21,310 --> 00:13:24,848
哎,鲁迪,别那么急,不要那么急嘛.

341
00:13:24,849 --> 00:13:30,886
好吧,夫人.那么我们在某处有这么一笔钱.

342
00:13:30,887 --> 00:13:33,056
我们将把它留给谁呢?

343
00:13:33,056 --> 00:13:34,085
好吧,

344
00:13:34,086 --> 00:13:35,985
我希望把所有的钱留给

345
00:13:35,985 --> 00:13:37,494
列沃伦.肯尼思.钱德勒.

346
00:13:37,495 --> 00:13:38,494
你知道这个人吗?

347
00:13:38,494 --> 00:13:40,723
他最近老是出现在

348
00:13:40,724 --> 00:13:41,993
达拉斯的电视里,

349
00:13:41,993 --> 00:13:44,332
他带着一头早熟的绻曲灰白头发,

350
00:13:44,332 --> 00:13:45,461
根本想不到

351
00:13:45,462 --> 00:13:46,831
去染一染,你知道吗?

352
00:13:46,831 --> 00:13:49,769
我希望他拿到这笔钱.

353
00:13:49,770 --> 00:13:53,098
对不起,柏蒂女士,嗯...

354
00:13:53,099 --> 00:13:54,868
什么?

355
00:13:54,869 --> 00:13:57,238
我只是有一个很实际的问题

356
00:13:57,238 --> 00:14:00,506
起草一个遗嘱或是任何别的契约

357
00:14:00,507 --> 00:14:03,775
要把家庭成员都剔除出去

358
00:14:03,775 --> 00:14:05,744
而取而代之的

359
00:14:05,745 --> 00:14:08,444
是转让这个遗产的一大部分

360
00:14:08,444 --> 00:14:13,282
给一个电视台的人物.

361
00:14:13,282 --> 00:14:15,781
啊,他是如同上帝一般的人.

362
00:14:15,781 --> 00:14:17,380
我感觉到了这一点.

363
00:14:17,381 --> 00:14:18,620
我懂.

364
00:14:18,620 --> 00:14:22,488
是不是还有别的途径使我们能够
处理这件事以更多的...

365
00:14:22,489 --> 00:14:24,258
你真的非把所有的一切都留给他吗?

366
00:14:24,258 --> 00:14:27,856
也许,比如说,比百分之二十五多一点儿?

367
00:14:27,857 --> 00:14:30,196
他要付很多的管理费,

368
00:14:30,196 --> 00:14:32,225
他的喷射飞机也越来越老.

369
00:14:32,225 --> 00:14:33,724
他的喷射飞机也越来越老?

370
00:14:33,725 --> 00:14:35,394
这样,鲁迪,我希望请你只是

371
00:14:35,394 --> 00:14:36,793
按照我要求的去起草

372
00:14:36,794 --> 00:14:37,893
然后拿回来让我过目

373
00:14:37,893 --> 00:14:40,602
以便审查它,怎么样?

374
00:14:40,602 --> 00:14:41,901
这就是那些个贼胚子

375
00:14:41,902 --> 00:14:44,231
在他们还是年轻和可爱的时候照的.

376
00:14:44,231 --> 00:14:48,229
剔除,剔除,剔除.

377
00:14:50,839 --> 00:14:52,238
你会回来的,是不是?

378
00:14:52,238 --> 00:14:54,737
哦,我会,我会的.

379
00:14:54,737 --> 00:14:56,376
谢谢您.

380
00:14:56,377 --> 00:14:58,876
不客气.谢谢你.

381
00:14:58,876 --> 00:15:02,144
后院是个小公寓吗?

382
00:15:02,145 --> 00:15:03,514
以前是的.

383
00:15:03,514 --> 00:15:04,913
你觉得我的花园怎么样?

384
00:15:04,914 --> 00:15:07,213
啊,这是个很漂亮的花园.

385
00:15:07,213 --> 00:15:08,882
您自己当园丁吗?

386
00:15:08,882 --> 00:15:10,321
多数都是自己做.

387
00:15:10,322 --> 00:15:12,751
我请了一个男孩子每周替我剪剪草地.

388
00:15:12,751 --> 00:15:14,820
付三十美金.你会相信吗?

389
00:15:14,820 --> 00:15:17,189
一般只要付五美金的.

390
00:15:17,189 --> 00:15:18,818
噢,我并不期望您有兴趣

391
00:15:18,819 --> 00:15:20,488
把这个地方出租,您想不想出租?

392
00:15:20,488 --> 00:15:22,327
我可以多付你租金,怎么样?

393
00:15:22,328 --> 00:15:25,256
我可以非常合理的租给你住

394
00:15:25,257 --> 00:15:27,366
假如你能在零碎活上帮我一把.

395
00:15:27,366 --> 00:15:28,865
那当然,那当然.

396
00:15:28,865 --> 00:15:31,733
也许从今以后拔拔杂草.

397
00:15:31,734 --> 00:15:34,732
当然,我是拔杂草专家.

398
00:15:36,802 --> 00:15:39,071
老是要往医院跑.

399
00:15:39,072 --> 00:15:41,741
伯鲁瑟和主要管区与他一起长大的家伙们

400
00:15:41,741 --> 00:15:42,910
有协议.

401
00:15:42,910 --> 00:15:45,139
他们每天早上都把事故报告传给他.

402
00:15:45,140 --> 00:15:46,609
我可以问你点儿事吗?

403
00:15:46,609 --> 00:15:47,678
当然.

404
00:15:47,679 --> 00:15:51,107
到底伯鲁瑟想让我做什么?

405
00:15:51,107 --> 00:15:52,916
去发掘案子,找到受害人.

406
00:15:52,917 --> 00:15:55,815
你让他们委托简.李曼.斯通律师事务所
替他们打官司.

407
00:15:55,816 --> 00:15:57,185
把案子集中起来.

408
00:15:57,185 --> 00:15:59,214
这么说我得施点儿手腕?

409
00:15:59,215 --> 00:16:01,254
你在法学院他们是怎么教你的?

410
00:16:01,254 --> 00:16:02,453
嗯,他们并没有教我追着救护车跑.

411
00:16:02,453 --> 00:16:05,721
好吧,你最好是快点儿学会这些,
否则你就要饿肚子了.

412
00:16:10,760 --> 00:16:12,989
你所要做的就是约一个门诊

413
00:16:12,990 --> 00:16:15,429
我们在这儿就可以为您登记.当然.

414
00:16:15,429 --> 00:16:16,628
好漂亮的花.

415
00:16:16,628 --> 00:16:18,027
多谢.

416
00:16:18,028 --> 00:16:19,157
维廉姆.

417
00:16:19,157 --> 00:16:20,166
狄克,你好吗?

418
00:16:20,167 --> 00:16:24,495
挺好.去346病房.

419
00:16:24,496 --> 00:16:25,535
威尔斯大夫.

420
00:16:25,535 --> 00:16:26,564
早上好.

421
00:16:26,565 --> 00:16:28,804
你好吗?有幸见到你.

422
00:16:28,804 --> 00:16:32,802
别整的像个律师.

423
00:16:36,611 --> 00:16:43,978
哼."请勿入内."

424
00:16:43,979 --> 00:16:46,248
你好吗?麦克肯兹先生?

425
00:16:46,248 --> 00:16:50,416
我的化验怎么样啦?

426
00:16:50,416 --> 00:16:54,384
胆囊气涨.钓错鱼了.

427
00:16:54,385 --> 00:16:56,824
万.蓝德尔先生.

428
00:16:56,824 --> 00:17:01,322
下午好,万.蓝德尔先生.

429
00:17:01,323 --> 00:17:04,421
能听见我说话吗,万.蓝德尔先生?

430
00:17:04,421 --> 00:17:05,590
嗨.

431
00:17:05,591 --> 00:17:07,030
你是谁啊?

432
00:17:07,031 --> 00:17:08,890
我是狄克.谢佛来尔,"解难"律师.

433
00:17:08,890 --> 00:17:10,029
你还没有和任何保险公司

434
00:17:10,030 --> 00:17:11,429
交谈过吧?有没有?

435
00:17:11,429 --> 00:17:12,428
没有.

436
00:17:12,429 --> 00:17:13,868
很好.不要和他们交谈

437
00:17:13,868 --> 00:17:14,997
因为他们只是要勒索你.

438
00:17:14,998 --> 00:17:16,027
你有没有私人律师?

439
00:17:16,027 --> 00:17:17,266
没有.

440
00:17:17,267 --> 00:17:19,736
我的律师事务所处理的车祸

441
00:17:19,736 --> 00:17:21,105
比孟菲斯任何其他人处理得都多.

442
00:17:21,106 --> 00:17:23,675
保险公司见了我们就怕,

443
00:17:23,675 --> 00:17:25,704
而我们却不收那怕一分钱.

444
00:17:25,704 --> 00:17:27,573
你能不能等我太太回来了再说?

445
00:17:27,573 --> 00:17:28,702
你太太,万先生...

446
00:17:28,703 --> 00:17:29,772
噢!

447
00:17:29,773 --> 00:17:32,112
啊!

448
00:17:32,112 --> 00:17:33,741
对不起.

449
00:17:33,741 --> 00:17:35,880
真对不起,万.蓝德尔先生.

450
00:17:35,880 --> 00:17:37,309
非常抱歉.

451
00:17:37,310 --> 00:17:39,849
你太太在那里,万.蓝德尔先生?

452
00:17:39,849 --> 00:17:42,078
她稍等一会儿就会回来.

453
00:17:42,078 --> 00:17:43,117
好的,我必须在我的办公室

454
00:17:43,118 --> 00:17:44,347
和她谈谈

455
00:17:44,347 --> 00:17:46,686
因为这里有大量的信息是我需要的.

456
00:17:46,687 --> 00:17:47,956
在这上面签个字吧.

457
00:17:47,956 --> 00:17:49,725
记住,除了你的大夫之外

458
00:17:49,725 --> 00:17:50,884
别与任何人交谈.

459
00:17:50,885 --> 00:17:52,424
会有各方面的人

460
00:17:52,425 --> 00:17:53,424
找上门来

461
00:17:53,424 --> 00:17:55,023
向你提供解决方案.

462
00:17:55,024 --> 00:17:57,523
我不希望你在任何情况下

463
00:17:57,523 --> 00:17:59,692
在我没有审查之前签署任何文件.

464
00:17:59,692 --> 00:18:01,561
听明白了没有?我的电话号码
就写在名片上.

465
00:18:01,561 --> 00:18:03,460
一天二十四小时你都可以找我.

466
00:18:03,461 --> 00:18:05,830
鲁迪.拜勒先生的电话号码印在反面.

467
00:18:05,830 --> 00:18:08,898
你也可以在任何时候找他.好吗?

468
00:18:08,899 --> 00:18:09,938
你还有什么疑问吗?

469
00:18:09,938 --> 00:18:11,037
没有.

470
00:18:11,038 --> 00:18:13,966
很好.我们会帮你讨回许多许多钱.

471
00:18:13,967 --> 00:18:15,606
我们走吧.

472
00:18:15,606 --> 00:18:17,375
对你的腿我真是非常抱歉.

473
00:18:17,376 --> 00:18:22,374
请你让我一个人清静一下.

474
00:18:24,513 --> 00:18:26,282
这就是怎么去做交易的.

475
00:18:26,283 --> 00:18:28,112
小菜一碟啦.

476
00:18:28,112 --> 00:18:30,681
要是遇到个家伙有自己的律师呢?

477
00:18:30,681 --> 00:18:32,780
我们空手而来,

478
00:18:32,780 --> 00:18:34,249
假使他别管是因为什么把我们踢出门.

479
00:18:34,250 --> 00:18:36,549
我们又失去了什呢?

480
00:18:36,549 --> 00:18:38,048
丢了点儿为人的尊严.

481
00:18:38,048 --> 00:18:40,487
也许是还有那么点儿君子的自重.

482
00:18:40,488 --> 00:18:43,187
你瞧...

483
00:18:43,187 --> 00:18:44,526
鲁迪,在法学院里

484
00:18:44,526 --> 00:18:46,085
他们并不教你你想学的.

485
00:18:46,086 --> 00:18:48,325
只不过是些定理啦,玄虚的概念啦

486
00:18:48,325 --> 00:18:50,594
再加上厚重的伦理学书本.

487
00:18:50,594 --> 00:18:51,993
伦理学又有什么不对?

488
00:18:51,993 --> 00:18:54,532
我猜,没什么不对吧.

489
00:18:54,533 --> 00:18:57,302
我的意思是,我相信一个律师

490
00:18:57,302 --> 00:18:58,531
必须为他的委托人去搏斗,

491
00:18:58,531 --> 00:18:59,670
防止钱财被盗去,

492
00:18:59,671 --> 00:19:00,830
并且尽量不欺骗.

493
00:19:00,830 --> 00:19:02,429
你知道的,这是最基本的.

494
00:19:02,430 --> 00:19:04,369
这就是去闹哄哄的追逐救护车.

495
00:19:04,369 --> 00:19:05,738
没错儿,可是谁去管这些?

496
00:19:05,739 --> 00:19:07,308
遍地都是律师.

497
00:19:07,308 --> 00:19:09,677
这就是交易市场.
这就是竞技场.

498
00:19:09,677 --> 00:19:13,145
他们没教你的是
在法学院你会被伤害了.

499
00:19:13,146 --> 00:19:16,944
你能判断出一个律师在行骗吗?

500
00:19:16,945 --> 00:19:19,214
他的嘴唇在动.

501
00:19:19,214 --> 00:19:24,082
什么又是一个妓女和一个律师的区别呢?

502
00:19:24,082 --> 00:19:27,880
一个妓女会在你死后停止勒索你.

503
00:19:27,881 --> 00:19:31,289
大家都喜欢听有关律师的笑话,

504
00:19:31,290 --> 00:19:32,489
特别是牵涉一群律师的.

505
00:19:32,489 --> 00:19:35,118
他们甚至以各式各样的这类笑话而自豪.

506
00:19:35,118 --> 00:19:38,326
你能从中得出什么结论吗?

507
00:19:38,327 --> 00:19:40,426
第三方并无什么区别.

508
00:19:40,426 --> 00:19:42,995
所以控方也可以利用这一点

509
00:19:42,995 --> 00:19:44,164
去引进第三方证人.

510
00:19:44,165 --> 00:19:47,333
你这是在干什么?

511
00:19:47,334 --> 00:19:50,063
噢,我在学习.

512
00:19:50,063 --> 00:19:53,501
我觉得我们要利用自己的时间去学习.

513
00:19:53,502 --> 00:19:55,031
我知道,伯鲁瑟,可是你瞧,先生,

514
00:19:55,031 --> 00:19:57,200
律师执照考试下个礼拜举行,对不对?

515
00:19:57,200 --> 00:19:58,769
就是下个礼拜.

516
00:19:58,770 --> 00:20:00,009
我有点儿担心,先生.

517
00:20:00,009 --> 00:20:01,738
嗨,鲁迪,瞧,你想学习,

518
00:20:01,739 --> 00:20:04,038
为什么你不能跑到医院去

519
00:20:04,038 --> 00:20:05,777
跟狄克学习?

520
00:20:05,777 --> 00:20:08,945
我不想从那方面学...

521
00:20:08,946 --> 00:20:10,415
我不想跟着狄克学.

522
00:20:10,416 --> 00:20:12,945
不错,我这里有个警察的报告.

523
00:20:12,945 --> 00:20:16,913
我们能代表这里的这个受害人吗?
雷科女士?

524
00:20:16,913 --> 00:20:17,982
嗯,还没有.

525
00:20:17,983 --> 00:20:19,412
那为什么你不能跑到医院去

526
00:20:19,413 --> 00:20:20,422
看看是怎么回事?

527
00:20:20,422 --> 00:20:25,980
也许你能得到她的委托.

528
00:20:59,808 --> 00:21:02,547
你能替我弄点饮料来吗?

529
00:21:02,547 --> 00:21:12,983
当然,亲爱的.

530
00:21:12,984 --> 00:21:19,021
你要的饮料来了.

531
00:21:19,021 --> 00:21:21,120
这些伤都是从哪儿来的?

532
00:21:21,121 --> 00:21:22,660
你想告诉他们

533
00:21:22,660 --> 00:21:25,289
这些伤是怎么来的吗?

534
00:21:25,289 --> 00:21:26,398
在我注视他们那一刻

535
00:21:26,399 --> 00:21:28,328
我知道这事是怎么发生的.

536
00:21:28,328 --> 00:21:30,467
就如同我十岁大时,

537
00:21:30,467 --> 00:21:31,766
我父亲在卧室里哭,

538
00:21:31,767 --> 00:21:33,136
我母亲满脸是血

539
00:21:33,136 --> 00:21:34,495
坐在餐桌旁

540
00:21:34,496 --> 00:21:35,565
告诉我说爸爸很懊悔,

541
00:21:35,566 --> 00:21:38,265
并保证以后再不这样做了.

542
00:21:38,265 --> 00:21:39,564
老天爷!

543
00:21:39,564 --> 00:21:43,572
你只要对我说是!

544
00:21:43,573 --> 00:21:45,572
为什么你要这样对我?

545
00:21:45,572 --> 00:21:49,140
和你在一起总是这样!

546
00:21:49,141 --> 00:21:55,008
你让我发疯!

547
00:21:55,009 --> 00:21:57,078
凯丽.雷科三天前

548
00:21:57,078 --> 00:21:59,277
就被送到圣.彼德医院了.

549
00:21:59,277 --> 00:22:01,246
我想强调一下,在半夜里,

550
00:22:01,247 --> 00:22:03,456
带着各式各样的伤.

551
00:22:03,456 --> 00:22:07,324
警察在小房间的沙发上发现了她,

552
00:22:07,324 --> 00:22:10,093
赤身裸体裹着条毯子,
几乎被打死了.

553
00:22:10,093 --> 00:22:11,092
克利夫.雷科,

554
00:22:11,093 --> 00:22:15,221
她的配偶,
明显是使用了毒品,

555
00:22:15,222 --> 00:22:16,391
处于高度骚动中,

556
00:22:16,391 --> 00:22:20,029
并在一开始把他带到他妻子面前时

557
00:22:20,030 --> 00:22:21,659
想打骂警察.

558
00:22:21,659 --> 00:22:22,828
而她,这么说吧,

559
00:22:22,829 --> 00:22:25,867
是被一根铝合金棒球棒狠狠地打过的.

560
00:22:25,868 --> 00:22:30,136
这就是他选择的凶器.

561
00:22:30,136 --> 00:22:33,074
让我们来讨论一下柏蒂女士的百万家产.

562
00:22:33,075 --> 00:22:36,403
不,我想讨论克利夫的案子.

563
00:22:36,404 --> 00:22:39,203
我想讨论克利夫身上发生了些什么.

564
00:22:39,203 --> 00:22:40,672
他被关押了一晚上.

565
00:22:40,673 --> 00:22:41,672
他家里把他保释了出来.

566
00:22:41,672 --> 00:22:43,281
他在一周内就要上堂.

567
00:22:43,282 --> 00:22:45,241
没什么事会发生.

568
00:22:45,241 --> 00:22:46,950
哦,不错.

569
00:22:46,950 --> 00:22:49,978
考林.简尼思.伯德松.

570
00:22:49,979 --> 00:22:52,148
她确凿从她后来第二个丈夫手里

571
00:22:52,149 --> 00:22:53,618
继承了近两百万遗产.

572
00:22:53,618 --> 00:22:57,956
可是律师们,无良的信托部门授权人,

573
00:22:57,956 --> 00:23:01,824
以及税务局的家伙们吞食了整个这笔财产.

574
00:23:01,825 --> 00:23:03,994
除了四十万美金之外的全部财产.

575
00:23:03,994 --> 00:23:07,862
柏蒂女士也许,
把它们塞进了床垫里,

576
00:23:07,863 --> 00:23:09,462
如同我们所知的那样.

577
00:23:09,462 --> 00:23:12,460
对不起.

578
00:23:18,699 --> 00:23:20,738
"你是我的阳光"

579
00:23:20,738 --> 00:23:22,037
"我唯一的阳光"

580
00:23:22,038 --> 00:23:24,607
"你令我快乐"

581
00:23:24,607 --> 00:23:26,536
"当天空灰濛濛的时候"

582
00:23:26,536 --> 00:23:29,534
"哒哒哒哒哒"

583
00:23:29,535 --> 00:23:34,243
"不要夺走我的阳光"

584
00:23:34,244 --> 00:23:35,813
噢,早上好啊,鲁迪,

585
00:23:35,813 --> 00:23:37,042
这是不是可爱的一天啊?

586
00:23:37,043 --> 00:23:40,351
噢,是的,天气真美.

587
00:23:40,351 --> 00:23:43,150
噢,我的护根材料到了.

588
00:23:43,150 --> 00:23:44,519
正好.那边点儿,靠那边点儿.

589
00:23:44,520 --> 00:23:46,319
就在这儿.对.

590
00:23:46,319 --> 00:23:48,148
停,停下.

591
00:23:48,149 --> 00:23:49,988
不错,就直接卸在这儿吧.

592
00:23:49,988 --> 00:23:51,687
我的园丁会来取.

593
00:23:51,687 --> 00:23:53,986
就卸在这儿.

594
00:23:53,987 --> 00:23:57,985
这是不是好大一堆护根材料啊?

595
00:24:24,376 --> 00:24:25,875
对不起.

596
00:24:25,875 --> 00:24:29,813
我并非多管闲事,但是...

597
00:24:29,814 --> 00:24:32,613
你没什么吧?

598
00:24:32,613 --> 00:24:35,052
你是不是痛啦?

599
00:24:35,052 --> 00:24:37,221
不是.

600
00:24:37,221 --> 00:24:39,580
不过还是多谢了.

601
00:24:39,581 --> 00:24:41,490
嗯,好吧,我就站这儿.

602
00:24:41,490 --> 00:24:43,319
我只不过在那边准备律师执照考试,

603
00:24:43,319 --> 00:24:45,958
要是你有什么事的话,

604
00:24:45,958 --> 00:24:47,487
就招呼我一声.

605
00:24:47,488 --> 00:24:48,787
好的.

606
00:24:48,787 --> 00:24:50,396
任何事情,好吗?

607
00:24:50,397 --> 00:24:55,795
我会为你去取.

608
00:24:55,795 --> 00:24:57,294
我叫鲁迪.拜勒.

609
00:24:57,294 --> 00:24:58,893
凯丽.雷科.很高兴与你相识.

610
00:24:58,894 --> 00:25:01,832
凯丽,能认识你很高兴.

611
00:25:01,833 --> 00:25:03,902
为什么你不坐下来?

612
00:25:03,902 --> 00:25:11,369
去搬个椅子坐下来嘛.

613
00:25:11,369 --> 00:25:16,807
你在那个学校念书?

614
00:25:16,807 --> 00:25:18,446
我进过奥斯汀.匹艾

615
00:25:18,447 --> 00:25:23,815
后来就进了孟菲斯州立大学的法学院.

616
00:25:23,815 --> 00:25:26,284
我一直想进大学,

617
00:25:26,284 --> 00:25:30,622
但却一直也没成功.

618
00:25:30,622 --> 00:25:31,751
真的吗?

619
00:25:31,752 --> 00:25:33,921
是啊,我一直以为我会进的.

620
00:25:33,921 --> 00:25:35,720
可就是没成功.

621
00:25:35,721 --> 00:25:38,789
你想成为一个什么样的律师?

622
00:25:38,790 --> 00:25:43,798
嗯,我喜欢干打官司的事,所以嘛,嗯...

623
00:25:43,798 --> 00:25:47,596
我喜欢整天都待在法庭上.

624
00:25:47,596 --> 00:25:50,964
罪犯的辩护律师吗?

625
00:25:50,965 --> 00:25:52,534
也许吧.

626
00:25:52,535 --> 00:25:53,964
也许是吧.

627
00:25:53,964 --> 00:25:57,002
他们会因一个出色的辩护而扬名.

628
00:25:57,003 --> 00:25:59,332
他们有权整天待在法院里.

629
00:25:59,332 --> 00:26:03,740
那杀人犯呢?

630
00:26:03,741 --> 00:26:09,948
大多数杀人犯雇不起私人律师.

631
00:26:09,949 --> 00:26:12,578
还有强奸犯以及...

632
00:26:12,578 --> 00:26:17,086
玩弄幼儿的的罪犯?

633
00:26:17,086 --> 00:26:18,785
不.

634
00:26:18,785 --> 00:26:23,783
那么殴打他们妻子的男人呢?

635
00:26:33,560 --> 00:26:36,798
处理犯罪的工作有其真正少有的特殊性.

636
00:26:36,799 --> 00:26:40,067
嗯,我也许要多做些,嗯...

637
00:26:40,068 --> 00:26:42,667
民事方面的诉讼.

638
00:26:43,697 --> 00:26:45,036
就像法律控告一类的事.

639
00:26:45,036 --> 00:26:47,994
--对.
--是的.

640
00:26:50,074 --> 00:26:53,702
那是,哦...请原谅.

641
00:27:00,481 --> 00:27:01,640
鲁迪.拜勒.

642
00:27:01,640 --> 00:27:03,709
你好.是我.

643
00:27:03,710 --> 00:27:04,979
你在干什么?

644
00:27:04,979 --> 00:27:06,948
钓鱼钓得怎么样?

645
00:27:06,948 --> 00:27:09,477
嗯,这个,嗯,进行的还不错.

646
00:27:09,478 --> 00:27:11,587
事实上我现在正在和
潜在的委托人交谈.

647
00:27:11,587 --> 00:27:13,956
不错.你最好使她签约雇你.

648
00:27:13,956 --> 00:27:16,984
怎么啦--鲁迪?鲁迪?
你听见我说吗?

649
00:27:16,985 --> 00:27:18,284
和你见面真令人愉快.

650
00:27:18,284 --> 00:27:19,583
鲁迪?

651
00:27:19,584 --> 00:27:21,823
是啊,多谢你和我作伴.

652
00:27:21,823 --> 00:27:24,262
嘿,明晚上怎么样?

653
00:27:24,262 --> 00:27:28,260
也许行.

654
00:27:41,436 --> 00:27:43,235
时间到了,搁笔吧.

655
00:27:43,236 --> 00:27:45,275
请把你的考卷往右传,

656
00:27:45,275 --> 00:27:47,174
这样监考员可以把它们集中起来.

657
00:27:47,174 --> 00:27:49,103
我在法学院的第一学年,

658
00:27:49,103 --> 00:27:51,272
大家都是相亲相爱的.

659
00:27:51,273 --> 00:27:53,342
因为我们都是学法律的,

660
00:27:53,342 --> 00:27:55,881
而法律是很清高的事.

661
00:27:55,881 --> 00:27:57,750
到了第三学年.

662
00:27:57,750 --> 00:28:00,549
你要庆幸自己没在睡梦里被人谋杀掉.

663
00:28:00,549 --> 00:28:01,748
人们在考场作弊,

664
00:28:01,749 --> 00:28:04,088
从图书馆里偷取研究资料,

665
00:28:04,088 --> 00:28:08,286
对教授们撒谎.

666
00:28:08,287 --> 00:28:13,285
这是这个行业的自然现象.

667
00:28:16,424 --> 00:28:17,923
就是这儿.

668
00:28:17,923 --> 00:28:20,432
半小时之前,她丈夫来到这里...

669
00:28:20,432 --> 00:28:22,391
把一大碗汤泼向她,

670
00:28:22,392 --> 00:28:24,401
因为她不想听那些

671
00:28:24,401 --> 00:28:26,130
"我是多么爱你"之类的甜言蜜语.

672
00:28:26,130 --> 00:28:28,369
这是我的病房.

673
00:28:28,370 --> 00:28:31,438
十八岁怀孕,结婚,然后就流产了--

674
00:28:31,439 --> 00:28:33,438
也许就因为他殴打她--

675
00:28:33,438 --> 00:28:35,537
但是,她却仍然不想离开他.

676
00:28:35,537 --> 00:28:37,876
你得帮我一下.

677
00:28:37,876 --> 00:28:39,205
我虽没亲眼看到这些

678
00:28:39,206 --> 00:28:41,605
却告诉了我这个姑娘是个失败者.

679
00:28:41,605 --> 00:28:44,513
流产啊断骨啊也许还有更危险的遭遇,

680
00:28:44,514 --> 00:28:47,382
以前真没见过任何人像她这样,

681
00:28:47,383 --> 00:28:48,612
与其逃避,

682
00:28:48,612 --> 00:28:53,610
我要尽我的一切保护她.

683
00:29:13,864 --> 00:29:18,702
探视时间已过.
孩子得被包起来了.

684
00:29:37,775 --> 00:29:39,874
那么...

685
00:29:39,874 --> 00:29:42,613
他该被枪毙.

686
00:29:42,613 --> 00:29:45,951
任何一个用铝合金棒球棒
殴打自己妻子的家伙

687
00:29:45,952 --> 00:29:48,880
都该被枪毙.

688
00:29:48,881 --> 00:29:50,890
你怎么会知道的?

689
00:29:50,890 --> 00:29:53,189
有警察的报告,
有救护车的报告,

690
00:29:53,190 --> 00:29:55,689
还有医院的病历.

691
00:29:55,689 --> 00:29:57,258
凯丽,你还想再等多久呢?

692
00:29:57,258 --> 00:29:59,357
一直到他决定用他的棒球棒

693
00:29:59,357 --> 00:30:01,026
打碎你的脑袋吗?

694
00:30:01,027 --> 00:30:03,226
这会致你于死地,
你知道吗,这绝对可能.

695
00:30:03,226 --> 00:30:04,665
对着脑壳来两颗描准的子弹,

696
00:30:04,666 --> 00:30:05,725
就是一样.

697
00:30:05,725 --> 00:30:14,871
别说了.
用不着告诉我我该怎么想.

698
00:30:14,872 --> 00:30:16,701
看着我,凯丽.

699
00:30:16,701 --> 00:30:22,139
我能问你点儿事吗?

700
00:30:22,139 --> 00:30:27,077
你有父亲或者兄弟吗?

701
00:30:27,078 --> 00:30:29,447
为什么问这个?

702
00:30:29,447 --> 00:30:30,776
因为假如你是我的女儿

703
00:30:30,776 --> 00:30:32,615
被人像你丈夫殴打你一样殴打了,

704
00:30:32,616 --> 00:30:37,284
我向上帝起誓我会扭断他的脖子.

705
00:30:37,284 --> 00:30:41,422
没有长兄吗?

706
00:30:41,423 --> 00:30:45,921
没有.

707
00:30:45,921 --> 00:30:53,598
没人照顾我.你知道吗?

708
00:30:53,598 --> 00:30:57,666
凯丽,我会尽我所能来帮助你,

709
00:30:57,667 --> 00:31:01,235
但你一定得申请离婚.

710
00:31:01,236 --> 00:31:04,134
趁还在医院里医治这最后一次
挨打的创伤时

711
00:31:04,135 --> 00:31:05,864
赶快办吧.

712
00:31:05,864 --> 00:31:07,373
它会使你从此扬眉吐气.

713
00:31:07,373 --> 00:31:10,241
还有什么事能证明是比这更好的吗?

714
00:31:10,242 --> 00:31:15,440
我...我不能申请离婚.

715
00:31:15,441 --> 00:31:17,540
为什么不能呢?

716
00:31:17,540 --> 00:31:21,748
因为他会杀了我.

717
00:31:21,748 --> 00:31:27,485
他一直是对我这么说的.

718
00:31:27,486 --> 00:31:29,085
这决不会发生.

719
00:31:29,086 --> 00:31:31,455
请你递个枕头给我,

720
00:31:31,455 --> 00:31:34,254
把它垫在我腿下好吗?

721
00:31:34,254 --> 00:31:45,430
那边有一只.

722
00:31:45,430 --> 00:31:49,428
这里.

723
00:31:49,429 --> 00:31:50,568
这样行吗?

724
00:31:50,568 --> 00:31:51,597
是.

725
00:31:51,598 --> 00:31:53,367
嗯,好吧.

726
00:31:53,367 --> 00:31:55,436
请你再把我的睡袍递给我好吗?

727
00:31:55,436 --> 00:31:58,005
好的.

728
00:31:58,005 --> 00:31:59,804
多谢.

729
00:31:59,805 --> 00:32:02,104
你需要我帮你穿上吗?

730
00:32:02,104 --> 00:32:04,373
不用.只需要背过身去就行.

731
00:32:04,373 --> 00:32:08,371
好吧.

732
00:32:25,586 --> 00:32:26,785
嗨!

733
00:32:26,785 --> 00:32:28,484
噢!你是谁?

734
00:32:28,485 --> 00:32:30,194
我住在这里.你又是什么人?

735
00:32:30,194 --> 00:32:32,123
噢,我的上帝.
我是笛尔伯特的太太.

736
00:32:32,123 --> 00:32:33,962
笛尔伯特?谁又是笛尔伯特?

737
00:32:33,963 --> 00:32:35,222
你怎么进来的?

738
00:32:35,222 --> 00:32:36,961
是柏蒂把钥匙交给我的.

739
00:32:36,962 --> 00:32:37,931
你是谁啊?

740
00:32:37,931 --> 00:32:39,290
我就是那个住在这里的人.

741
00:32:39,291 --> 00:32:40,300
你搞明白了没有?

742
00:32:40,301 --> 00:32:41,400
你无权待在这里.

743
00:32:41,400 --> 00:32:42,629
这是私人住宅.

744
00:32:42,630 --> 00:32:44,399
噢,是的,对了.
某个地方,抓牢些个.

745
00:32:44,399 --> 00:32:48,397
柏蒂要见你.

746
00:32:49,837 --> 00:32:51,406
"协议书是永久性的,

747
00:32:51,407 --> 00:32:54,935
且上述限制性包括在其中..."

748
00:32:54,935 --> 00:32:56,604
这是什么?

749
00:32:56,605 --> 00:32:58,604
嗯,你想必就是那个律师了.

750
00:32:58,604 --> 00:32:59,643
我是鲁迪.拜勒.

751
00:32:59,644 --> 00:33:02,243
我是笛尔伯特.伯德松.
柏蒂的小儿子.

752
00:33:02,243 --> 00:33:03,542
他朝我咋咋呼呼的.

753
00:33:03,542 --> 00:33:05,811
他让我滚出他的公寓.

754
00:33:05,811 --> 00:33:06,780
是这样吗?

755
00:33:06,781 --> 00:33:08,050
你说对了.
就是这样.

756
00:33:08,051 --> 00:33:09,080
对你们两位都是如此.

757
00:33:09,080 --> 00:33:10,579
我不希望你们中任何一个上那里去

758
00:33:10,580 --> 00:33:11,619
乱翻我的东西.

759
00:33:11,619 --> 00:33:12,618
那是私人财产.

760
00:33:12,619 --> 00:33:13,948
我回家来看望我的妈妈,

761
00:33:13,949 --> 00:33:16,917
发觉,他妈的,她让一个混身发臭的
律师和她住在一起.

762
00:33:16,918 --> 00:33:19,816
是你把我妈妈的遗嘱搞得一团糟?

763
00:33:19,817 --> 00:33:24,125
哼,她是你妈妈.
为什么你不直接问她?

764
00:33:24,125 --> 00:33:25,824
她一句话都不说.

765
00:33:25,824 --> 00:33:27,023
噢,好嘛,

766
00:33:27,024 --> 00:33:29,693
那我也不能讲任何话.

767
00:33:29,693 --> 00:33:37,600
这是律师与委托人间的特权.

768
00:33:37,600 --> 00:33:42,098
但是让我来告诉你.

769
00:33:42,099 --> 00:33:44,438
我打过好几个电话...

770
00:33:44,438 --> 00:33:45,467
嗯-哼.

771
00:33:45,467 --> 00:33:47,236
查过某些账目.

772
00:33:47,237 --> 00:33:49,776
你妈妈的第二任丈夫

773
00:33:49,776 --> 00:33:52,205
给她留下了巨额的财产.

774
00:33:52,205 --> 00:33:53,614
多大呢?

775
00:33:53,615 --> 00:33:57,243
非常巨大.

776
00:33:57,243 --> 00:34:00,751
我希望你没有插手,小伙子,

777
00:34:00,752 --> 00:34:04,450
妈妈,你是不是爱到佛罗利达我们家去

778
00:34:04,451 --> 00:34:06,320
看看我们并住上一段时间?

779
00:34:06,320 --> 00:34:08,289
妈妈,你会爱上那儿的.

780
00:34:08,289 --> 00:34:10,318
妈妈,坐回来吃点儿点心吧,妈妈.

781
00:34:10,319 --> 00:34:11,588
我来给她倒杯牛奶.

782
00:34:11,588 --> 00:34:13,987
我们住的离青年喷水池很近.

783
00:34:13,987 --> 00:34:19,984
我们虽不在近旁但也只离开三十三英里远
而离开迪斯尼世界不过一百八十英里.

784
00:34:42,137 --> 00:34:44,146
丹尼.雷,感觉怎么样?

785
00:34:44,147 --> 00:34:45,276
非常好.

786
00:34:45,276 --> 00:34:47,915
我看上去不是很棒吗?

787
00:34:47,915 --> 00:34:49,144
你觉得能行吗?

788
00:34:49,145 --> 00:34:50,884
行啊.动身吧.

789
00:34:50,884 --> 00:34:54,352
好吧.

790
00:34:54,353 --> 00:34:56,752
律师是本不应该把自己搅进

791
00:34:56,752 --> 00:34:58,481
他的委托人的日常生活中的,

792
00:34:58,482 --> 00:35:01,091
但是有各种各样的律师

793
00:35:01,091 --> 00:35:06,089
也同样有着各种各样的委托人.

794
00:35:17,835 --> 00:35:31,340
都还好吧,宝贝儿?

795
00:35:31,340 --> 00:35:34,079
雷科!雷科!雷科!

796
00:35:34,079 --> 00:35:35,908
雷科!雷科!

797
00:35:35,908 --> 00:35:38,547
还需要点儿什么别的吗?

798
00:35:38,547 --> 00:35:41,775
嗯,还想要点儿别的吗?

799
00:35:41,776 --> 00:35:46,744
请再给我一杯"杰克.丹尼尔"吧.

800
00:35:46,745 --> 00:35:51,953
干嘛?想宰了我啊?

801
00:35:51,953 --> 00:35:54,652
柏蒂女士,这是我的朋友丹尼.雷.

802
00:35:54,652 --> 00:35:55,751
噢,丹尼.雷,

803
00:35:55,751 --> 00:35:58,759
欢迎你的光临.

804
00:35:58,760 --> 00:36:00,329
您好,柏蒂女士.

805
00:36:00,330 --> 00:36:02,089
在--噢.

806
00:36:02,089 --> 00:36:03,328
就这儿.

807
00:36:03,329 --> 00:36:05,498
你就过来坐在这里吧.

808
00:36:05,498 --> 00:36:08,726
放松一下.

809
00:36:08,727 --> 00:36:10,266
您快把我吃光了,

810
00:36:10,266 --> 00:36:11,965
在这里挡一下.

811
00:36:11,966 --> 00:36:16,964
嗯...

812
00:36:21,402 --> 00:36:23,401
喔...

813
00:36:23,402 --> 00:36:26,980
对不起,柏蒂女士,我得歇会儿,

814
00:36:26,980 --> 00:36:28,339
腰酸背痛啊.

815
00:36:28,340 --> 00:36:36,247
我差点儿忘了,这是寄给你的.

816
00:36:36,247 --> 00:36:37,586
啊,我的上帝.

817
00:36:37,587 --> 00:36:39,646
呕!太让人激动了不是?

818
00:36:39,646 --> 00:36:42,754
呕.我真为你感到骄傲.

819
00:36:42,755 --> 00:36:43,854
我通过律师执照考试了.

820
00:36:43,854 --> 00:36:44,853
太棒了.

821
00:36:44,854 --> 00:36:47,892
嗨,来来,这杯敬鲁迪.

822
00:36:47,893 --> 00:36:50,891
祝贺你通过考试.

823
00:36:50,892 --> 00:36:52,661
那是什么玩意儿?

824
00:36:52,661 --> 00:36:55,460
冰茶.

825
00:36:55,460 --> 00:36:56,729
祝贺你,鲁迪.

826
00:36:56,730 --> 00:36:59,329
孩子们,干得不错啊.

827
00:36:59,329 --> 00:37:01,428
我今天从处理万.蓝德尔一事的决定里
取到一张支票.

828
00:37:01,428 --> 00:37:02,927
我决定发给你们奖金.

829
00:37:02,928 --> 00:37:05,467
每人五千五百美金.

830
00:37:05,467 --> 00:37:06,566
谢谢你.

831
00:37:06,566 --> 00:37:09,804
再争取多签些合同,啊?

832
00:37:09,805 --> 00:37:10,934
好吧.

833
00:37:10,935 --> 00:37:12,074
嗯.

834
00:37:12,074 --> 00:37:15,972
是不是为明天一早准备就绪啦?

835
00:37:15,973 --> 00:37:19,611
九点整?

836
00:37:19,612 --> 00:37:21,481
我们要争论

837
00:37:21,481 --> 00:37:24,749
巨大福利保险公司的驳回提案.

838
00:37:24,750 --> 00:37:28,018
嗯,对,我觉得我们准备好了.
一切就绪.

839
00:37:28,019 --> 00:37:31,887
狄克和我已经滤过一遍.
我想我们准备好了.

840
00:37:31,887 --> 00:37:34,955
但愿如此,因为我可能,嗯...

841
00:37:34,956 --> 00:37:40,124
我可能让你出马参加一些辩论,鲁迪.

842
00:37:40,124 --> 00:37:44,492
要是我们因为案子驳回而
输了这场官司

843
00:37:44,493 --> 00:37:46,122
那将是非常令人苦恼的.

844
00:37:46,122 --> 00:37:48,361
老板?

845
00:37:48,362 --> 00:37:51,530
好吧,我必须走了.

846
00:37:51,530 --> 00:37:53,099
我来付账.

847
00:37:53,100 --> 00:37:55,729
你们享用吧.

848
00:37:55,729 --> 00:38:00,727
多谢,伯鲁瑟,非常感谢.

849
00:38:04,506 --> 00:38:06,205
有些情况不妙.

850
00:38:06,205 --> 00:38:08,704
我非常确信这一点.

851
00:38:08,704 --> 00:38:12,172
他以前从来不这么分钱.

852
00:38:12,173 --> 00:38:15,481
你能想像在他的慷慨背后是什么吗?

853
00:38:15,482 --> 00:38:18,550
得了吧,伙计.

854
00:38:18,551 --> 00:38:21,619
噢,是吗?

855
00:38:21,620 --> 00:38:24,289
瞧这里.

856
00:38:24,289 --> 00:38:26,818
昨天,伯鲁瑟过去的一个搭档

857
00:38:26,818 --> 00:38:28,157
对大法官作证.

858
00:38:28,157 --> 00:38:29,986
我觉得他和大法官有了协议.

859
00:38:29,987 --> 00:38:31,456
也许只是个时间问题,

860
00:38:31,456 --> 00:38:32,085
他就会对着伯鲁瑟高唱凯歌.

861
00:38:32,086 --> 00:38:34,095
那又怎么样?

862
00:38:34,095 --> 00:38:37,063
怎么样?到时后你就得拍屁股走人.

863
00:38:37,064 --> 00:38:40,022
走人?狄克,我才刚进来.

864
00:38:40,023 --> 00:38:44,861
是吗?嗯,事情也许会稍微热闹起来.

865
00:38:44,862 --> 00:38:46,861
贿赂陪审员,逃税,

866
00:38:46,861 --> 00:38:48,700
敲竹杠,应有尽有.

867
00:38:48,700 --> 00:38:50,399
我是害怕,真的害怕.

868
00:38:50,400 --> 00:38:53,398
你怕什么呢?他们又不能抓你我.

869
00:38:53,399 --> 00:38:55,668
听着...

870
00:38:55,668 --> 00:38:58,676
假设他们带着传票和锯子来呢.嗯?

871
00:38:58,677 --> 00:39:00,006
他们可以这么干.

872
00:39:00,006 --> 00:39:02,974
他们在处理诈骗犯的案子里这么干过.

873
00:39:02,975 --> 00:39:04,774
他们过来,没收所有文件,

874
00:39:04,775 --> 00:39:07,184
拿走所有的计算机.
你手里还剩什么?

875
00:39:07,184 --> 00:39:08,983
我并不担心被抓起来.

876
00:39:08,983 --> 00:39:10,652
我担心我的饭碗.

877
00:39:10,652 --> 00:39:13,311
好吧,那你说了这些是什么意思?

878
00:39:13,312 --> 00:39:15,951
开溜.

879
00:39:15,951 --> 00:39:18,490
你有多少钱?

880
00:39:18,490 --> 00:39:21,458
我...我有五千五百美金.

881
00:39:21,459 --> 00:39:22,818
我也一样.

882
00:39:22,818 --> 00:39:24,927
我们能租一个小办公室,

883
00:39:24,927 --> 00:39:26,586
五百美金一月吧.

884
00:39:26,587 --> 00:39:30,765
我们小本经营它六个月.

885
00:39:30,765 --> 00:39:32,694
它将会漂亮起来.

886
00:39:32,695 --> 00:39:34,834
我们平分一切,所有的一切.

887
00:39:34,834 --> 00:39:37,303
一半一半.

888
00:39:37,303 --> 00:39:40,161
开销,花费,利润--一切的一切

889
00:39:40,162 --> 00:39:44,670
从中对劈.

890
00:39:44,670 --> 00:39:47,469
怎么啦?

891
00:39:47,469 --> 00:39:49,608
你不愿意和我合伙吗?

892
00:39:49,609 --> 00:39:51,538
啊,不是,不是那意思--他们--

893
00:39:51,538 --> 00:39:53,877
不是那意思.只是你...

894
00:39:53,877 --> 00:39:55,776
你刚才真是给了我一击,

895
00:39:55,777 --> 00:39:58,146
我的意思是,你得给我点儿时间,行不行?

896
00:39:58,146 --> 00:39:59,775
我是说,你不能就靠吓唬我一下就完了.

897
00:39:59,775 --> 00:40:05,253
我们必须快速行动.

898
00:40:05,253 --> 00:40:07,922
好吧.那就动手吧.

899
00:40:07,922 --> 00:40:10,850
我们做一阵子,

900
00:40:10,851 --> 00:40:12,350
看看行不行.

901
00:40:12,351 --> 00:40:14,590
你手里有几个案子?

902
00:40:14,590 --> 00:40:17,389
嗯,我...我有三个.

903
00:40:17,389 --> 00:40:20,058
把它们取出来,放到家里去.

904
00:40:20,058 --> 00:40:24,426
可别让人抓住,明白?

905
00:40:24,426 --> 00:40:26,065
有人监视我们吗?

906
00:40:26,066 --> 00:40:29,164
联邦政府.

907
00:40:29,165 --> 00:40:33,163
我吃地太快了.

908
00:41:00,923 --> 00:41:02,622
对不起.

909
00:41:02,623 --> 00:41:06,551
嗨,听着,今天早上我去了办公室.

910
00:41:06,551 --> 00:41:08,890
联邦调查局把大门给封了.

911
00:41:08,891 --> 00:41:10,060
伯鲁瑟到了吗?

912
00:41:10,060 --> 00:41:11,859
我怀疑这一点.

913
00:41:11,860 --> 00:41:13,229
已经发出了对伯鲁瑟和

914
00:41:13,229 --> 00:41:14,358
普林斯的逮捕令.

915
00:41:14,359 --> 00:41:15,598
噢,家伙.

916
00:41:15,598 --> 00:41:17,627
放松点儿.你能担起来的.

917
00:41:17,628 --> 00:41:19,767
这只是个提案.进行吧.

918
00:41:19,767 --> 00:41:21,136
-我?
-对了.来吧.

919
00:41:21,136 --> 00:41:22,235
我要撑不住了.

920
00:41:22,236 --> 00:41:23,735
你熟悉这个案子.
你会处理好的.

921
00:41:23,735 --> 00:41:25,134
听我说,我还没领到执照呢.

922
00:41:25,135 --> 00:41:28,933
我们用不着执照.
走吧.

923
00:41:35,011 --> 00:41:37,480
我绝对记住了勃拉克一案的详情.

924
00:41:37,481 --> 00:41:39,580
我读了法学中所有的书,

925
00:41:39,580 --> 00:41:42,009
及关于证据和证物的规定,

926
00:41:42,009 --> 00:41:44,078
但当我瞥了一眼法庭的四周,

927
00:41:44,078 --> 00:41:48,886
却觉得自己像个完全无助的婴儿.

928
00:41:48,887 --> 00:41:53,285
原谅我,哦,先生,
尊敬的法官,哦,呵呵.

929
00:41:53,285 --> 00:41:57,523
我来参加勃拉克状告
巨大福利保险公司的听证.

930
00:41:57,523 --> 00:41:59,052
那么你是谁?

931
00:41:59,053 --> 00:42:00,352
哦,鲁迪.拜勒.

932
00:42:00,352 --> 00:42:03,760
我为简.李曼.斯通律师事务所工作.

933
00:42:03,761 --> 00:42:09,958
噢,你为简.李曼...工作.

934
00:42:09,959 --> 00:42:15,437
斯涛切是烟草游说团的支持者.

935
00:42:15,437 --> 00:42:17,606
斯通先生在哪儿呢?

936
00:42:17,606 --> 00:42:20,065
哦...

937
00:42:20,065 --> 00:42:22,404
老实讲,尊敬的法官,
我...我不知道.

938
00:42:22,405 --> 00:42:24,444
他本来应该在这里和我汇合,

939
00:42:24,444 --> 00:42:26,443
但我不知道他现在在哪里.

940
00:42:26,443 --> 00:42:29,212
嗯,这怎么会不令我吃惊呢?

941
00:42:29,212 --> 00:42:32,210
那么你想干什么?
你想继续下去.

942
00:42:32,211 --> 00:42:33,380
不是,尊敬的法官.

943
00:42:33,931 --> 00:42:37,059
我...我准备来争辩那个提案.

944
00:42:37,059 --> 00:42:39,798
你是个律师吗?

945
00:42:39,798 --> 00:42:42,137
嗯,我刚通过了律师执照考试,

946
00:42:42,138 --> 00:42:44,267
而,哦,这些人是我的委托人.

947
00:42:44,267 --> 00:42:46,336
斯通先生只是代表我立案

948
00:42:46,336 --> 00:42:47,835
以等我通过执照考试.

949
00:42:47,836 --> 00:42:50,105
哼,鬼知道你怎么会如此胆大包天

950
00:42:50,105 --> 00:42:52,074
没有执照就敢迈进我的法庭.

951
00:42:52,074 --> 00:42:53,903
立即从这里滚出去,
去取得你的执照,

952
00:42:53,904 --> 00:42:55,313
然后再回到这儿来.

953
00:42:55,313 --> 00:42:56,842
去拿个执照!

954
00:42:56,842 --> 00:42:58,911
到你有了该死的执照再回来!

955
00:42:58,912 --> 00:43:02,110
要是这能令法院满意...

956
00:43:02,111 --> 00:43:05,419
根据记录,我的名字是莱奥.富.楚门

957
00:43:05,419 --> 00:43:07,248
来自霆雷-博瑞特律师事务所,

958
00:43:07,249 --> 00:43:08,718
巨大福利保险公司的法律顾问,

959
00:43:08,718 --> 00:43:10,217
我要说

960
00:43:10,218 --> 00:43:13,786
要是这个年轻人
已经通过了律师执照考试,

961
00:43:13,786 --> 00:43:16,555
尊敬的法官,就让他
参加本案的辩论吧.

962
00:43:16,555 --> 00:43:20,653
我们欢迎他参加,嗯,
重大时刻的诉讼.

963
00:43:20,654 --> 00:43:23,323
你没有异义吗,楚门先生?

964
00:43:23,323 --> 00:43:25,362
假如这能令法院满意,

965
00:43:25,362 --> 00:43:28,061
我愿意很荣幸的来介绍拜勒先生

966
00:43:28,061 --> 00:43:29,460
为大田纳西州的

967
00:43:29,461 --> 00:43:31,230
法律实习律师.

968
00:43:31,230 --> 00:43:34,068
法官,您是可以现在就让他宣誓的,

969
00:43:34,069 --> 00:43:36,997
而我也将很高兴支持他.

970
00:43:36,998 --> 00:43:40,106
你肯定你准备好了要
进行下去吗,拜勒先生?

971
00:43:40,107 --> 00:43:42,975
绝对--

972
00:43:42,976 --> 00:43:45,775
是的.尊敬的法官,是这样.

973
00:43:45,775 --> 00:43:49,113
很好,那么,举起你的右手.

974
00:43:49,114 --> 00:43:51,043
你是否庄严地宣誓,鲁迪.拜勒,

975
00:43:51,043 --> 00:43:53,542
你将忠实地,正直地维护美国

976
00:43:53,542 --> 00:43:55,311
和田纳西州的

977
00:43:55,312 --> 00:43:56,651
宪法和法律,

978
00:43:56,651 --> 00:43:58,720
在实践中引导你自己

979
00:43:58,721 --> 00:44:03,519
去最大限度地运用你的技巧和能力,
来帮助你的上帝?

980
00:44:03,519 --> 00:44:05,058
我发誓.

981
00:44:05,058 --> 00:44:08,556
嗯,很好,很好.这就是了.
让我们进行下去吧.

982
00:44:08,557 --> 00:44:09,926
-恭喜你.
-谢谢你.

983
00:44:09,927 --> 00:44:11,256
欢迎投入战斗.

984
00:44:11,256 --> 00:44:14,994
由一个白痴来施行就职宣誓
而又是由一个恶棍来做担保.

985
00:44:14,995 --> 00:44:15,994
哈维.

986
00:44:15,994 --> 00:44:17,723
我终于是个律师了.

987
00:44:17,724 --> 00:44:23,401
请进.

988
00:44:23,744 --> 00:44:25,443
请允许我,哈维?

989
00:44:25,443 --> 00:44:27,912
当然了.

990
00:44:27,912 --> 00:44:31,750
嗯,请坐.

991
00:44:31,751 --> 00:44:34,609
嗯...

992
00:44:34,610 --> 00:44:39,218
这场官司烦扰着我,拜勒先生.

993
00:44:39,218 --> 00:44:41,487
我不想用"琐碎的"这个词,

994
00:44:41,487 --> 00:44:44,256
而坦白地说

995
00:44:44,256 --> 00:44:45,685
其是非曲直也没给我什么印象.

996
00:44:45,686 --> 00:44:48,754
"嗯呵.我想法在锻炼时间溜出来一下."

997
00:44:48,755 --> 00:44:51,124
事实上,我真是厌倦了这类官司.

998
00:44:51,124 --> 00:44:54,022
"别说没借口."

999
00:44:54,023 --> 00:44:57,061
我倾向于答应驳回的提案.

1000
00:44:57,062 --> 00:45:02,400
现在,你可以,嗯,重新向
联邦法院立案,你知道的,

1001
00:45:02,400 --> 00:45:05,099
到别的什么地方去打这个官司.

1002
00:45:05,099 --> 00:45:08,097
我不想让它塞满我的备审目录.

1003
00:45:08,098 --> 00:45:10,467
原谅我,我得去一下厕所.

1004
00:45:10,467 --> 00:45:15,675
"你也一样."

1005
00:45:15,675 --> 00:45:19,043
鲁迪...

1006
00:45:19,044 --> 00:45:21,583
我是个很贵的律师

1007
00:45:21,583 --> 00:45:23,942
属于一个很贵的律师事务所,嗯...

1008
00:45:23,942 --> 00:45:25,911
并且我手里有很多案子要办.

1009
00:45:25,912 --> 00:45:28,521
我可以投掷飞标非常接近目标的中心.

1010
00:45:28,521 --> 00:45:31,549
我告诉我的委托人
巨大福利保险公司,

1011
00:45:31,550 --> 00:45:35,018
打这个官司要花很多的钱.

1012
00:45:35,019 --> 00:45:37,887
而对于你和你的委托人,也一样.

1013
00:45:37,888 --> 00:45:41,256
这样,他们授权我...

1014
00:45:41,256 --> 00:45:47,363
许给你和你的委托人...

1015
00:45:47,364 --> 00:45:50,632
五万美金来和解这桩案子.

1016
00:45:50,633 --> 00:45:54,631
而我...我甚至可以加码

1017
00:45:54,632 --> 00:45:59,770
这个数,哦,到七万五千美金.

1018
00:45:59,770 --> 00:46:02,908
不再承担责任,你明白的.

1019
00:46:02,909 --> 00:46:05,438
当然,嗯...

1020
00:46:05,438 --> 00:46:06,877
要是你拒绝,

1021
00:46:06,877 --> 00:46:12,245
那么就会引发第三次世界大战.

1022
00:46:12,245 --> 00:46:12,744
该我了,能吗,哈唯?

1023
00:46:12,745 --> 00:46:16,713
当然啦.

1024
00:46:18,513 --> 00:46:21,212
哦...

1025
00:46:22,352 --> 00:46:24,681
我看这里没有多少

1026
00:46:24,681 --> 00:46:26,220
官司好打,我的孩子.

1027
00:46:26,220 --> 00:46:29,558
可也许我能依靠莱奥

1028
00:46:29,559 --> 00:46:31,858
为和解出个价.

1029
00:46:31,858 --> 00:46:33,827
他们也许愿意付些钱给你

1030
00:46:33,828 --> 00:46:37,826
以避免一分钟付莱奥好几千块.

1031
00:46:37,826 --> 00:46:42,194
然而他已经为辩方向我出了价了.

1032
00:46:42,195 --> 00:46:45,163
噢,真的吗?多少?

1033
00:46:45,164 --> 00:46:47,933
哦,七万五.

1034
00:46:47,933 --> 00:46:50,702
我的天.

1035
00:46:50,702 --> 00:46:52,871
嗯,仔细想想,我的孩子,

1036
00:46:52,871 --> 00:46:57,269
要是你不接受的话
你一定是疯了.

1037
00:46:57,269 --> 00:46:58,708
您这么认为吗?

1038
00:46:58,709 --> 00:47:01,378
七万五?天哪.

1039
00:47:01,378 --> 00:47:03,647
这是...
这可是笔大钱哪.

1040
00:47:03,647 --> 00:47:05,116
这可真不像是莱奥干的.

1041
00:47:05,117 --> 00:47:08,815
嗯,我也这么想.
他是个大人物.

1042
00:47:08,815 --> 00:47:12,053
嗯.

1043
00:47:12,054 --> 00:47:16,052
呼.

1044
00:47:16,053 --> 00:47:17,082
怎么样?

1045
00:47:17,082 --> 00:47:18,551
在房间里没有开会.

1046
00:47:18,552 --> 00:47:20,091
那是一场伏击.

1047
00:47:20,091 --> 00:47:21,050
发生了什么事?

1048
00:47:21,051 --> 00:47:22,490
他们插标把我卖了.

1049
00:47:22,491 --> 00:47:23,590
是啊.

1050
00:47:23,590 --> 00:47:25,559
你以为他们敢这么对待伯鲁瑟吗?

1051
00:47:25,559 --> 00:47:27,088
不会.对这些伯鲁瑟是太了解了.

1052
00:47:27,089 --> 00:47:28,558
--对.
--那么怎么回事?

1053
00:47:28,558 --> 00:47:29,757
他们许给我七万五.

1054
00:47:29,758 --> 00:47:31,257
好啊,接过来.

1055
00:47:31,257 --> 00:47:33,626
我们应得的三分之一就是两万五.

1056
00:47:33,627 --> 00:47:35,096
我们正需要钱.

1057
00:47:35,096 --> 00:47:37,095
海尔法官对驳回这个案子
是认真的.

1058
00:47:37,095 --> 00:47:38,564
他只是个又老又暴燥的人,

1059
00:47:38,565 --> 00:47:40,064
让他坐在那个位置上是太常了些.

1060
00:47:40,064 --> 00:47:42,533
我说,我们所能做的最好情况
也就是能得多少得多少

1061
00:47:42,533 --> 00:47:45,371
以减轻他父母的负担.

1062
00:47:45,372 --> 00:47:47,901
保险公司出了价以求和解.

1063
00:47:47,902 --> 00:47:50,111
开的什么价?

1064
00:47:50,111 --> 00:47:53,709
七万五千美金.

1065
00:47:53,709 --> 00:47:55,808
他们算出那是为打官司

1066
00:47:55,809 --> 00:47:57,848
雇律师所要花的钱.

1067
00:47:57,848 --> 00:47:59,117
噢.

1068
00:47:59,118 --> 00:48:03,346
狗娘样的以为他们能收买我们.

1069
00:48:03,346 --> 00:48:09,223
他们确实这么以为.

1070
00:48:09,224 --> 00:48:15,491
你最好和他谈谈.

1071
00:48:15,492 --> 00:48:20,790
你要和解吗,鲁迪?

1072
00:48:20,790 --> 00:48:27,667
我是说,其中一部份钱是归你的.

1073
00:48:27,667 --> 00:48:31,395
决不.

1074
00:48:31,396 --> 00:48:35,834
我想揭露这些人.

1075
00:48:35,835 --> 00:48:38,843
妈妈,能请给我一杯水吗?

1076
00:48:38,844 --> 00:48:48,180
当然,宝贝.

1077
00:48:48,180 --> 00:48:50,149
不管从这个官司拿回什么,

1078
00:48:50,150 --> 00:48:55,148
用它来照顾我的家.

1079
00:48:55,148 --> 00:48:57,347
我真的爱他们.

1080
00:48:57,347 --> 00:48:59,356
我止不住要想

1081
00:48:59,356 --> 00:49:01,855
当我们大家都不相信会死去那样

1082
00:49:01,855 --> 00:49:03,884
在边缘的四周挣扎,

1083
00:49:03,885 --> 00:49:05,854
这个孩子却直视着它,

1084
00:49:05,854 --> 00:49:08,393
无助地看着它.

1085
00:49:08,393 --> 00:49:10,192
简直不能设想

1086
00:49:10,192 --> 00:49:14,190
这要有多大的勇气.

1087
00:49:19,399 --> 00:49:20,528
喂.

1088
00:49:20,529 --> 00:49:22,898
猜猜昨晚上谁死了.

1089
00:49:22,898 --> 00:49:23,897
谁?

1090
00:49:23,898 --> 00:49:25,407
你根本没睡过吧?

1091
00:49:25,407 --> 00:49:26,866
哈维.海尔.

1092
00:49:26,867 --> 00:49:27,866
62岁.

1093
00:49:27,866 --> 00:49:29,875
豪门贵胄啊.

1094
00:49:29,875 --> 00:49:31,404
海尔法官?

1095
00:49:31,405 --> 00:49:33,374
对.心脏病突发.

1096
00:49:33,374 --> 00:49:35,413
跌进他的游泳池死了.

1097
00:49:35,414 --> 00:49:36,913
你是在开玩笑吧.

1098
00:49:36,913 --> 00:49:37,912
嗯哼.

1099
00:49:37,913 --> 00:49:39,382
你一定是开玩笑吧.

1100
00:49:39,382 --> 00:49:40,381
嗯嗯.

1101
00:49:40,382 --> 00:49:42,411
猜猜新上任的法官

1102
00:49:42,411 --> 00:49:44,410
如何处理巨大福利保险公司案?

1103
00:49:44,410 --> 00:49:46,419
见鬼这我怎么会知道?

1104
00:49:46,420 --> 00:49:47,619
泰荣,开普勒,

1105
00:49:47,619 --> 00:49:49,918
哈佛毕业的黑人民权律师.

1106
00:49:49,918 --> 00:49:51,917
他不会支持霆雷-博瑞特律师事务所的,

1107
00:49:51,918 --> 00:49:53,887
他严厉对待那些保险公司--

1108
00:49:53,887 --> 00:49:55,426
一直都在控诉他们.

1109
00:49:55,426 --> 00:49:56,885
对于我们确是大幸.

1110
00:49:56,886 --> 00:49:58,095
慢着,慢着.

1111
00:49:58,095 --> 00:49:59,894
小孩儿,你知道什么是造雨人吗?

1112
00:49:59,895 --> 00:50:01,424
是钱从天上掉下来啊!

1113
00:50:01,424 --> 00:50:02,893
现在是五点,你什么时候到?

1114
00:50:02,894 --> 00:50:04,393
立马就到.

1115
00:50:04,393 --> 00:50:06,432
赶快来.我们要讨论一下
有关快速立法的事.

1116
00:50:06,433 --> 00:50:07,392
咖啡马上就到.

1117
00:50:07,392 --> 00:50:10,260
好的.

1118
00:50:10,261 --> 00:50:12,900
楚门先生...

1119
00:50:12,900 --> 00:50:14,069
你对这个案子的快速立法

1120
00:50:14,070 --> 00:50:15,899
提出异义.

1121
00:50:15,899 --> 00:50:17,898
理由是什么?

1122
00:50:17,898 --> 00:50:19,907
嗯,尊敬的法官,哦...

1123
00:50:19,908 --> 00:50:21,937
这个案子已被

1124
00:50:21,937 --> 00:50:23,936
海尔大法官裁决了.

1125
00:50:23,936 --> 00:50:26,904
哦,这个案子要快速立法

1126
00:50:26,905 --> 00:50:28,414
所需的准备工作

1127
00:50:28,415 --> 00:50:30,344
我相信对双方都是

1128
00:50:30,344 --> 00:50:31,313
过度的负担.

1129
00:50:31,314 --> 00:50:33,643
毫无意义.

1130
00:50:33,643 --> 00:50:36,511
我想问问你,楚门先生.

1131
00:50:36,512 --> 00:50:38,011
作为一个辩护律师,

1132
00:50:38,011 --> 00:50:39,920
你是否曾经赞成过

1133
00:50:39,921 --> 00:50:42,320
法律诉讼的快速立法?

1134
00:50:42,320 --> 00:50:44,449
嗯,尊敬的法官,
我想信我赞成过.

1135
00:50:45,419 --> 00:50:46,888
好吧.

1136
00:50:46,888 --> 00:50:51,586
那么告诉我那个案子的名字
以及处理它的法庭.

1137
00:50:51,587 --> 00:50:52,926
嗯,呵呵...

1138
00:50:52,926 --> 00:50:53,925
哦,尊敬的法官,

1139
00:50:53,926 --> 00:50:56,225
我得回去找一下才能告诉你.

1140
00:50:56,225 --> 00:51:00,193
好吧,下午三点给我个电话.

1141
00:51:00,194 --> 00:51:03,332
我想三点钟之前我到不了办公室.

1142
00:51:03,332 --> 00:51:05,931
噢.好吧,你到了办公室之后就给我个电话.

1143
00:51:05,932 --> 00:51:08,401
我非常渴望能听到你赞同快速立法

1144
00:51:08,401 --> 00:51:09,940
的那个案子的情况.

1145
00:51:09,940 --> 00:51:11,439
是,先生.

1146
00:51:11,440 --> 00:51:14,438
那个孩子正面临死亡,先生们.

1147
00:51:14,438 --> 00:51:17,476
你们同意我们需要为他的证词录音.

1148
00:51:17,477 --> 00:51:18,946
是,确实是这样.

1149
00:51:18,947 --> 00:51:20,446
当然,尊敬的法官.

1150
00:51:20,446 --> 00:51:22,475
只是在我的出庭备忘录

1151
00:51:22,476 --> 00:51:23,945
压得我喘不过气来.

1152
00:51:23,945 --> 00:51:27,183
下个礼拜四下午怎么样?

1153
00:51:27,184 --> 00:51:30,852
对我来说很好,尊敬的法官.

1154
00:51:30,853 --> 00:51:32,352
原谅我,尊敬的法官.

1155
00:51:32,352 --> 00:51:33,951
距今天有整一个礼拜.

1156
00:51:33,952 --> 00:51:35,951
我相信我会出城.

1157
00:51:35,951 --> 00:51:39,989
是的,礼拜四我得出城.

1158
00:51:39,989 --> 00:51:43,627
那么作证就定在...

1159
00:51:43,628 --> 00:51:46,157
下个礼拜四下午两点.

1160
00:51:46,157 --> 00:51:49,625
很遗憾假如这对辩方带来不方便,

1161
00:51:49,626 --> 00:51:54,664
但是上帝了解你们有足够的人
来处理这件事.

1162
00:51:54,664 --> 00:51:55,903
好了,下一个问题是什么?

1163
00:51:55,904 --> 00:51:57,403
噢,哦,尊敬的法官,

1164
00:51:57,403 --> 00:52:02,371
驳回此案的提案还是悬而未决的.

1165
00:52:02,372 --> 00:52:03,871
啊,对了.

1166
00:52:03,871 --> 00:52:06,839
那个提案被否决了.

1167
00:52:06,840 --> 00:52:08,339
嗯,

1168
00:52:08,339 --> 00:52:10,238
我想这事就是这样了.

1169
00:52:10,239 --> 00:52:13,147
先生们,我们走吧.

1170
00:52:13,148 --> 00:52:14,647
对您的任命表示祝贺,

1171
00:52:14,647 --> 00:52:15,646
尊敬的法官.

1172
00:52:15,647 --> 00:52:19,385
谢谢你,先生.

1173
00:52:19,386 --> 00:52:23,384
楚门先生...

1174
00:52:23,384 --> 00:52:24,913
别忘了给我电话

1175
00:52:24,914 --> 00:52:28,992
告诉我你赞成快速立法一案的名字.

1176
00:52:28,992 --> 00:52:39,398
等我找一找.

1177
00:52:39,398 --> 00:52:44,866
你是不是觉得难于理解.我的孩子?

1178
00:52:44,867 --> 00:52:47,865
绝对是.

1179
00:52:49,905 --> 00:52:53,373
目前狄克想在地方上的饶舌音乐电台做个广告.

1180
00:52:53,374 --> 00:52:55,373
虽然我一直非常爱听饶舌音乐,

1181
00:52:55,373 --> 00:52:56,572
我们却负担不起广告费.

1182
00:52:56,572 --> 00:52:59,371
家具,法庭费,七百五买了个传真机,

1183
00:52:59,371 --> 00:53:00,910
四百租了台计算机,

1184
00:53:00,911 --> 00:53:03,410
付新买的二手车的首期,

1185
00:53:03,410 --> 00:53:05,479
我们再次破产了.

1186
00:53:05,479 --> 00:53:10,477
狄克说是他对这种情况感到毛骨悚然.

1187
00:53:17,595 --> 00:53:20,593
(凯丽.雷科)

1188
00:53:35,329 --> 00:53:37,328
喔.

1189
00:53:39,567 --> 00:53:40,896
鲁迪,嗨.

1190
00:53:40,897 --> 00:53:43,106
嗨,达特.

1191
00:53:43,106 --> 00:53:44,565
请进.你好吗?

1192
00:53:44,565 --> 00:53:46,074
我很好.你好吗?

1193
00:53:46,075 --> 00:53:47,304
很好.

1194
00:53:47,304 --> 00:53:49,903
听我说,礼拜三晚上我要到克利夫兰去.

1195
00:53:49,904 --> 00:53:51,443
那里是巨大福利保险公司总部所在地.

1196
00:53:51,443 --> 00:53:53,912
我要取得所有执行官的证词.

1197
00:53:53,912 --> 00:53:54,911
噢.

1198
00:53:54,912 --> 00:53:56,241
所以...

1199
00:53:56,241 --> 00:53:58,410
也用不着担心开销.

1200
00:53:58,411 --> 00:53:59,950
我们会处理这些事的.

1201
00:53:59,950 --> 00:54:00,949
谢谢你,鲁迪.

1202
00:54:00,950 --> 00:54:04,218
没问题.

1203
00:54:05,588 --> 00:54:08,117
真不好用.

1204
00:54:10,286 --> 00:54:11,615
我的合伙人狄克.谢佛来尔.

1205
00:54:11,616 --> 00:54:12,855
你好吗?

1206
00:54:12,855 --> 00:54:16,423
嗨.

1207
00:54:25,631 --> 00:54:27,460
嗨,鲁迪...

1208
00:54:27,460 --> 00:54:29,429
泰荣来了.

1209
00:54:29,430 --> 00:54:30,969
好的.

1210
00:54:30,969 --> 00:54:32,098
大法官到了.

1211
00:54:32,099 --> 00:54:33,268
噢.

1212
00:54:33,268 --> 00:54:34,937
哎,等一下,达特,把这给我.

1213
00:54:34,938 --> 00:54:35,967
噢.

1214
00:54:35,967 --> 00:54:39,305
谢谢你.

1215
00:54:39,306 --> 00:54:40,305
您好.

1216
00:54:40,306 --> 00:54:41,635
你好.

1217
00:54:41,635 --> 00:54:43,644
请进.

1218
00:54:43,645 --> 00:54:44,604
谢谢你.

1219
00:54:44,604 --> 00:54:46,113
这是勃拉克女士.

1220
00:54:46,114 --> 00:54:47,713
他就是大法官泰荣,开普勒.

1221
00:54:47,713 --> 00:54:49,282
非常荣幸见到您.

1222
00:54:49,283 --> 00:54:50,412
应该是我感到荣幸.

1223
00:54:50,412 --> 00:54:53,450
噢,哦,请过来.

1224
00:54:53,451 --> 00:54:55,820
这地方太小.

1225
00:54:55,820 --> 00:54:56,879
挤得很.

1226
00:54:56,880 --> 00:54:59,119
我们到室外去看看.

1227
00:54:59,119 --> 00:55:02,117
这么做对你来说行吗,我的孩子?

1228
00:55:02,118 --> 00:55:04,117
当然,可以的.

1229
00:55:04,117 --> 00:55:05,116
好的.

1230
00:55:05,117 --> 00:55:06,126
很好.

1231
00:55:06,127 --> 00:55:10,125
我来带路.

1232
00:55:20,232 --> 00:55:22,231
我看着每小时赚一千美金的

1233
00:55:22,231 --> 00:55:24,270
律师班子,

1234
00:55:24,270 --> 00:55:25,639
我恨他们.

1235
00:55:25,640 --> 00:55:27,739
从他们高高的栖息处,

1236
00:55:27,739 --> 00:55:30,238
朝整个司法系统

1237
00:55:30,238 --> 00:55:32,007
肆无忌惮地撒尿.

1238
00:55:32,007 --> 00:55:35,245
我从前憎恨他们
是因为我不是他们的对手.

1239
00:55:35,246 --> 00:55:38,005
而现在我恨他们
是因为他们在代表谁

1240
00:55:38,005 --> 00:55:40,074
和代表着的事.

1241
00:55:40,075 --> 00:55:42,014
尊敬的法官,您好吗?

1242
00:55:42,014 --> 00:55:43,083
很好,先生.

1243
00:55:43,083 --> 00:55:45,082
但愿这些狗没有吓着你.

1244
00:55:45,083 --> 00:55:46,752
我们要在室外做这件事了.

1245
00:55:46,752 --> 00:55:48,251
室内稍嫌拥挤了些.

1246
00:55:48,252 --> 00:55:49,781
找个位子坐下吧.

1247
00:55:49,781 --> 00:55:50,580
好吧.

1248
00:55:50,581 --> 00:55:55,229
--你好吗?
--很好.

1249
00:56:09,794 --> 00:56:12,792
谢谢你.

1250
00:56:18,431 --> 00:56:20,900
你好吗,丹尼.雷?

1251
00:56:20,900 --> 00:56:22,269
你好.

1252
00:56:22,270 --> 00:56:25,808
你已经见过大法官开普勒了.

1253
00:56:25,808 --> 00:56:27,277
是的.

1254
00:56:27,278 --> 00:56:30,346
这是莱奥.富.楚门

1255
00:56:30,347 --> 00:56:33,945
和他的同事们.

1256
00:56:33,945 --> 00:56:36,444
这是谭密.她是法院的记录员.

1257
00:56:36,445 --> 00:56:38,444
嗨.

1258
00:56:38,444 --> 00:56:39,953
让他起誓吧.

1259
00:56:39,953 --> 00:56:41,282
你是否发誓

1260
00:56:41,283 --> 00:56:42,822
你将给出的证词

1261
00:56:42,822 --> 00:56:44,321
是事实的,全都是事实,

1262
00:56:44,322 --> 00:56:46,091
除了事实没有其他,来帮助你的上帝?

1263
00:56:46,091 --> 00:56:49,489
是,我起誓.

1264
00:56:53,229 --> 00:56:54,728
我告诉过他了...

1265
00:56:54,728 --> 00:56:56,297
我知道,我知道.过来,过来,过来.

1266
00:56:56,298 --> 00:56:58,297
别担心,别担心,这只是个听证.

1267
00:56:58,297 --> 00:57:00,526
丹尼.雷,我是莱奥.楚门,

1268
00:57:00,526 --> 00:57:02,035
我代表巨大福利保险公司,

1269
00:57:02,036 --> 00:57:04,165
我非常非常遗憾...

1270
00:57:04,165 --> 00:57:05,804
我非常遗憾

1271
00:57:05,804 --> 00:57:07,363
在这种情形之下来到这里.

1272
00:57:07,364 --> 00:57:08,873
嗯,要是你的委托人做了他们

1273
00:57:08,873 --> 00:57:10,802
本该做的事,我们将不会聚于此地.

1274
00:57:10,802 --> 00:57:12,671
对不起,我没听清楚.

1275
00:57:12,672 --> 00:57:13,671
嘿,孩子,

1276
00:57:13,671 --> 00:57:15,310
想要块口香糖吗?

1277
00:57:15,311 --> 00:57:16,310
当然了.

1278
00:57:16,311 --> 00:57:17,840
给你.你胳膊断了吗?

1279
00:57:17,840 --> 00:57:19,179
对啦.

1280
00:57:19,180 --> 00:57:20,679
你遭遇意外了吗?

1281
00:57:20,679 --> 00:57:21,638
干嘛问哪?

1282
00:57:21,639 --> 00:57:23,808
嗯,我是个律师,并且,嗯...

1283
00:57:23,808 --> 00:57:26,317
把这交给你的母亲--
您是他母亲吗?

1284
00:57:26,317 --> 00:57:30,015
--嗯.
--对于这种意外事故我也许能给你们找点儿钱回来.

1285
00:57:42,861 --> 00:57:45,160
帮忙做点儿展示好吗?

1286
00:57:45,160 --> 00:57:47,829
好的,我先帮一下这位先生.

1287
00:57:47,829 --> 00:57:51,827
噢,不不,我来.这样挺好.

1288
00:57:54,997 --> 00:57:56,836
我能为您做点儿什么,先生?

1289
00:57:56,836 --> 00:58:01,384
噢,我只是看看.

1290
00:58:20,378 --> 00:58:22,847
街上有座电影院.

1291
00:58:22,847 --> 00:58:24,846
买张后排,中间的票.

1292
00:58:24,846 --> 00:58:27,854
我三十分钟后到.

1293
00:58:27,855 --> 00:58:30,853
好的.

1294
00:59:07,041 --> 00:59:13,218
克利夫想让我生个孩子.

1295
00:59:13,219 --> 00:59:17,047
嗯,你得做个决定.

1296
00:59:17,048 --> 00:59:19,087
他对性事着了魔.

1297
00:59:19,087 --> 00:59:21,416
他认为这样就能把我们栓在一起.

1298
00:59:21,416 --> 00:59:24,924
听我说,我真的不想讨论这些.

1299
00:59:24,925 --> 00:59:28,423
我只是想见见你.

1300
01:00:13,398 --> 01:00:15,397
我相信到克利夫兰去的一路上

1301
01:00:15,397 --> 01:00:17,436
我都能闻到她的香水味.

1302
01:00:17,436 --> 01:00:21,664
这使我难于集中精力去想
莱奥.楚门和他的那帮人.

1303
01:00:21,665 --> 01:00:23,234
他们会乘头等舱飞过来.

1304
01:00:23,234 --> 01:00:24,933
在悠闲的聚餐之后,

1305
01:00:24,934 --> 01:00:26,903
在个会议室里碰头

1306
01:00:26,903 --> 01:00:28,472
讨论怎样把我整个的毁掉.

1307
01:00:28,473 --> 01:00:30,512
当我得在
汽车旅馆定个房间的时候,

1308
01:00:30,512 --> 01:00:32,611
他们却正在高级套房里呼呼大睡.

1309
01:00:32,611 --> 01:00:35,639
他们养精蓄锐醒过来之后

1310
01:00:35,640 --> 01:00:37,609
可以整装待发投入战斗.

1311
01:00:37,609 --> 01:00:39,648
这是我要的听证会,

1312
01:00:39,649 --> 01:00:42,987
却在他们的地盘上举行.

1313
01:00:42,987 --> 01:00:44,486
啊,年轻的鲁迪.拜勒.

1314
01:00:44,487 --> 01:00:45,956
刚好到时间.

1315
01:00:45,956 --> 01:00:46,955
楚门先生.

1316
01:00:46,956 --> 01:00:49,015
泰勒,给这个孩子弄杯咖啡.

1317
01:00:49,015 --> 01:00:50,054
贾克.安德豪尔.

1318
01:00:50,055 --> 01:00:51,084
鲁迪.拜勒.

1319
01:00:51,085 --> 01:00:52,324
年轻律师永远是饥渴的.

1320
01:00:52,324 --> 01:00:54,823
所有这些个孩子都请求,
书面的...

1321
01:00:54,823 --> 01:00:56,692
一定有上百年了法律的经验

1322
01:00:56,693 --> 01:00:58,322
是在长桌周围收集起来的.

1323
01:00:58,322 --> 01:01:01,860
我的工作人员是个
六次律师执照考试失败者.

1324
01:01:01,861 --> 01:01:03,330
噢,鲁迪,别让

1325
01:01:03,330 --> 01:01:05,359
桌子这边这么多人给吓着.

1326
01:01:05,360 --> 01:01:07,329
我保证,要是你在高尔夫球场
碰到他们,

1327
01:01:07,329 --> 01:01:10,867
他们就像团成一团的廉价衣服.

1328
01:01:10,868 --> 01:01:11,837
让我瞧瞧.

1329
01:01:11,837 --> 01:01:14,206
让我瞧瞧我们要干些什么.

1330
01:01:14,206 --> 01:01:17,134
我想也许是,哦,从公司

1331
01:01:17,135 --> 01:01:20,173
的设计师贾克.安德豪尔

1332
01:01:20,174 --> 01:01:21,643
这里开始是合适的.

1333
01:01:21,644 --> 01:01:24,842
我...我并不这样认为.

1334
01:01:24,843 --> 01:01:26,182
对不起?

1335
01:01:26,182 --> 01:01:27,281
嗯,你听到我说的是什么.

1336
01:01:27,282 --> 01:01:28,811
我想从,哦,贾奇.莱曼奇科开始,

1337
01:01:28,811 --> 01:01:32,809
她是索赔处理人员.

1338
01:01:36,489 --> 01:01:40,157
我觉得还是最好从安德豪尔先生开始.

1339
01:01:40,157 --> 01:01:41,686
以所有可给予的尊敬,
楚门先生,

1340
01:01:41,687 --> 01:01:43,186
这是我要的听证会.

1341
01:01:43,186 --> 01:01:45,855
我得按我认为合式的顺序
召见这些证人,

1342
01:01:45,855 --> 01:01:50,863
所以我愿意从贾奇.莱曼奇科开始.

1343
01:01:50,863 --> 01:01:55,661
也许我们要打个电话给大法官.

1344
01:01:55,662 --> 01:01:58,660
哦,我不认为这么一大早

1345
01:01:58,661 --> 01:02:00,270
我们就非得运用拳击术.

1346
01:02:00,270 --> 01:02:03,598
我没有注册成为一个拳击师.

1347
01:02:03,599 --> 01:02:05,668
我们只是有那么点儿小问题

1348
01:02:05,668 --> 01:02:07,707
在贾奇.莱曼奇科

1349
01:02:07,707 --> 01:02:09,606
哦,这位波兰女士身上.

1350
01:02:09,607 --> 01:02:13,805
是什么样的问题呢?

1351
01:02:13,805 --> 01:02:16,773
她已经不再在这里工作了.

1352
01:02:16,774 --> 01:02:18,073
她被开除了吗?

1353
01:02:18,074 --> 01:02:22,482
她辞职了.

1354
01:02:22,482 --> 01:02:25,151
嗯,那么,她现在在哪儿?

1355
01:02:25,151 --> 01:02:26,520
嗯...

1356
01:02:26,521 --> 01:02:29,989
她不再为我们的委托人工作,所以,嗯,

1357
01:02:29,990 --> 01:02:32,719
我们不能再把她作为证人,

1358
01:02:32,719 --> 01:02:38,357
让我们继续下去.

1359
01:02:38,357 --> 01:02:40,356
好吧.鲁塞尔.克劳凯特.

1360
01:02:40,356 --> 01:02:45,194
在座的谁是鲁塞尔.克劳凯特?

1361
01:02:45,194 --> 01:02:47,263
他也离开了.

1362
01:02:47,263 --> 01:02:49,232
机构精简.

1363
01:02:49,233 --> 01:02:50,502
机构精简?

1364
01:02:50,502 --> 01:02:52,731
嗯,真会这么巧.

1365
01:02:52,732 --> 01:02:55,201
我们的委托人进行了
阶段性的机构精简.

1366
01:02:55,201 --> 01:02:59,209
是啊,这是会发生的,不是吗?

1367
01:03:05,207 --> 01:03:07,706
那么艾沃特,哦,鲁富肯...

1368
01:03:07,706 --> 01:03:09,245
那个索赔部付总裁呢?

1369
01:03:09,246 --> 01:03:12,074
他也被精简了吗?

1370
01:03:12,075 --> 01:03:16,073
不,他就在这里.

1371
01:03:17,713 --> 01:03:19,222
你是艾沃特.鲁富肯?

1372
01:03:19,222 --> 01:03:21,821
嗯哼.

1373
01:03:21,821 --> 01:03:23,720
我祝贺你...

1374
01:03:23,721 --> 01:03:25,250
鲁富肯先生...

1375
01:03:25,250 --> 01:03:27,759
在后来巨大福利保险公司的屠杀中

1376
01:03:27,759 --> 01:03:31,997
奇迹般的幸存下来.

1377
01:03:32,997 --> 01:03:37,495
呼哇.

1378
01:03:37,496 --> 01:03:42,494
我也许不敢说百分之百的有准备,
但我精神尚好.

1379
01:03:47,822 --> 01:03:49,541
我很好奇

1380
01:03:49,621 --> 01:03:51,310
好奇什么？

1381
01:03:52,220 --> 01:03:54,709
我很想知道...

1382
01:03:55,924 --> 01:03:58,264
你还记得自己第一次出卖自己是什么时候吗？

1383
01:04:04,675 --> 01:04:08,075
你小子很狂啊？

1384
01:04:11,176 --> 01:04:13,406
我建议你注意一下你的态度。

1385
01:04:13,476 --> 01:04:14,966
你可是律师

1386
01:04:15,046 --> 01:04:16,636
你在乎我的态度？

1387
01:04:18,917 --> 01:04:21,477
我这么远从孟菲斯来

1388
01:04:21,557 --> 01:04:23,487
来为4个人取证

1389
01:04:23,558 --> 01:04:25,998
居然有两个不在

1390
01:04:26,058 --> 01:04:29,158
而你让我注意态度？

1391
01:04:29,228 --> 01:04:31,258
这个，是你运气不好罢了

1392
01:04:31,339 --> 01:04:33,239
你打算怎么样啊，小子？

1393
01:04:35,509 --> 01:04:36,909
我要从他那里取证

1394
01:04:36,970 --> 01:04:39,240
就是坐那里的那位拉佛金先生

1395
01:04:39,310 --> 01:04:40,470
然后我会收拾行李

1396
01:04:40,540 --> 01:04:42,640
回孟菲斯去

1397
01:04:46,351 --> 01:04:49,251
这本来是没有保险的人才会碰到的

1398
01:04:49,321 --> 01:04:51,751
在这个医学如此进步的社会里

1399
01:04:51,822 --> 01:04:53,762
有着那么多优秀的医生

1400
01:04:53,822 --> 01:04:57,452
让这个孩子这么死掉简直是无耻之极

1401
01:04:57,533 --> 01:04:59,433
他有保单在身，

1402
01:04:59,503 --> 01:05:01,803
她母亲也付了保险金

1403
01:05:01,873 --> 01:05:04,803
也许不是很多的钱，但是足够了

1404
01:05:04,874 --> 01:05:06,774
我在这个案子中孤立无援

1405
01:05:06,844 --> 01:05:10,244
我被人耍，我也很怕

1406
01:05:10,315 --> 01:05:11,935
但是我是正确的

1407
01:05:12,015 --> 01:05:15,315
我坐在这个可怜的孩子旁边，看着他受苦

1408
01:05:15,385 --> 01:05:17,615
我发誓我一定要报复

1409
01:05:29,397 --> 01:05:30,297
喂？

1410
01:05:30,368 --> 01:05:31,798
鲁迪，是我

1411
01:05:31,868 --> 01:05:33,338
怎么了？

1412
01:05:34,378 --> 01:05:35,838
求你帮帮我

1413
01:05:35,908 --> 01:05:37,808
你在哪里？

1414
01:05:37,879 --> 01:05:39,309
店里

1415
01:05:39,379 --> 01:05:41,309
好，待那里别动，好吗？

1416
01:05:41,379 --> 01:05:42,849
我马上就来

1417
01:05:42,919 --> 01:05:44,039
好的

1418
01:05:59,532 --> 01:06:00,592
我是鲁迪.贝勒

1419
01:06:04,443 --> 01:06:06,173
她在这里

1420
01:06:11,954 --> 01:06:14,014
谢谢你能来这

1421
01:06:24,966 --> 01:06:26,626
发生什么了？

1422
01:06:26,696 --> 01:06:28,366
快点，我们给他看看

1423
01:06:44,218 --> 01:06:46,418
我扶着你呢，别怕

1424
01:06:46,489 --> 01:06:47,459
喔...

1425
01:06:47,519 --> 01:06:49,889
鲁迪，我真高兴你能打电话来

1426
01:06:49,959 --> 01:06:51,389
快进来

1427
01:06:51,460 --> 01:06:53,400
可怜的孩子

1428
01:06:53,460 --> 01:06:56,430
我会照顾你的

1429
01:06:56,500 --> 01:06:57,560
别怕

1430
01:06:57,640 --> 01:06:59,600
别怕，你会没事的

1431
01:06:59,671 --> 01:07:00,731
好的

1432
01:07:00,801 --> 01:07:02,431
你知道哪儿能找到我吧

1433
01:07:02,511 --> 01:07:03,911
知道，知道

1434
01:07:03,971 --> 01:07:05,441
这时我对自己说

1435
01:07:05,512 --> 01:07:07,442
我一定要尽我所能

1436
01:07:07,512 --> 01:07:08,912
帮她离婚

1437
01:07:08,982 --> 01:07:11,612
不然这家伙肯定会打死她为止

1438
01:07:11,683 --> 01:07:12,913
肯定的

1439
01:07:34,006 --> 01:07:35,606
这...

1440
01:07:36,676 --> 01:07:39,446
他本来死不了的

1441
01:07:39,517 --> 01:07:41,107
混蛋

1442
01:07:41,187 --> 01:07:42,447
迭克...

1443
01:07:49,058 --> 01:07:50,958
多特，我真难过

1444
01:08:40,716 --> 01:08:43,056
真谢谢你能来

1445
01:08:43,126 --> 01:08:44,616
非常感谢

1446
01:09:14,461 --> 01:09:16,221
会没事的

1447
01:09:21,902 --> 01:09:23,892
嘿!

1448
01:09:28,303 --> 01:09:29,693
这是布奇，他来找窃听器

1449
01:09:31,143 --> 01:09:32,073
窃听器？

1450
01:09:32,143 --> 01:09:33,233
嘘～嘘～

1451
01:09:44,965 --> 01:09:47,425
办公室内有小型窃听器

1452
01:09:55,207 --> 01:09:56,467
谢谢你

1453
01:09:56,537 --> 01:09:57,807
谢谢你

1454
01:09:57,877 --> 01:09:59,277
请用餐

1455
01:10:00,678 --> 01:10:03,118
这个窃听器

1456
01:10:03,178 --> 01:10:05,448
有中级电路

1457
01:10:07,219 --> 01:10:09,279
可以传送很微弱的信号

1458
01:10:10,519 --> 01:10:13,619
可能是捷克制的

1459
01:10:14,830 --> 01:10:16,630
不，我不认为...

1460
01:10:16,700 --> 01:10:19,290
警察或者联邦调查局会用这种类型的

1461
01:10:20,531 --> 01:10:22,631
肯定是其他什么人在窃听

1462
01:10:22,701 --> 01:10:25,171
还有谁想窃听我们？

1463
01:10:25,241 --> 01:10:26,141
嗯

1464
01:10:26,211 --> 01:10:28,301
我想到了

1465
01:10:32,152 --> 01:10:33,082
鲁迪

1466
01:10:33,153 --> 01:10:34,313
怎么了？

1467
01:10:34,383 --> 01:10:35,643
我只是随便打个电话

1468
01:10:35,723 --> 01:10:37,153
你要我从城里带什么吗？

1469
01:10:37,223 --> 01:10:38,683
不用

1470
01:10:38,753 --> 01:10:40,693
嘿，猜猜谁打算庭外和解了？

1471
01:10:40,754 --> 01:10:41,854
谁？

1472
01:10:41,924 --> 01:10:43,364
多特.布莱克

1473
01:10:43,424 --> 01:10:45,054
多特.布莱克？

1474
01:10:45,134 --> 01:10:46,684
我今天过去了一下

1475
01:10:46,765 --> 01:10:48,665
顺便看了看她

1476
01:10:48,735 --> 01:10:50,165
我给她买了个水果蛋糕

1477
01:10:50,235 --> 01:10:53,175
她说她已经没有劲头

1478
01:10:53,235 --> 01:10:55,705
打这种长官司了

1479
01:10:55,776 --> 01:10:56,676
还有...

1480
01:10:56,746 --> 01:10:57,676
她要多少？

1481
01:10:57,746 --> 01:10:59,606
她说17万5就行

1482
01:10:59,676 --> 01:11:01,376
我觉得我们该拿

1483
01:11:01,447 --> 01:11:03,717
先这么着吧，明天见

1484
01:11:03,787 --> 01:11:05,547
好吧，不过我建议你接受这个结果

1485
01:11:05,617 --> 01:11:07,347
我知道，我听见你的话了，我想想看

1486
01:11:07,418 --> 01:11:08,618
好吧，再见

1487
01:11:09,588 --> 01:11:10,648
哈哈

1488
01:11:11,798 --> 01:11:13,058
孩子...

1489
01:11:13,128 --> 01:11:16,358
这个家庭已经受了够多的苦了

1490
01:11:16,429 --> 01:11:19,269
我觉得那女人会接受庭外和解的

1491
01:11:20,169 --> 01:11:21,799
我会问问她的

1492
01:11:23,710 --> 01:11:25,900
你现在就给她打电话吧

1493
01:11:25,980 --> 01:11:28,170
我就在这等着

1494
01:11:28,251 --> 01:11:30,371
好，我一会打给你，德拉门德先生

1495
01:11:30,451 --> 01:11:31,541
再见

1496
01:11:42,963 --> 01:11:45,393
我对这男孩的死很难过

1497
01:11:45,463 --> 01:11:46,903
嗯

1498
01:11:46,963 --> 01:11:48,993
听着，我的客户，呃...

1499
01:11:49,074 --> 01:11:50,934
想要和解，鲁迪

1500
01:11:51,004 --> 01:11:52,404
我们这么说吧，鲁迪

1501
01:11:52,474 --> 01:11:54,904
只是讨论一下数目而已

1502
01:11:54,975 --> 01:11:56,445
假设这个保险被赔付了

1503
01:11:56,515 --> 01:11:59,415
我的客户大概要付...

1504
01:11:59,485 --> 01:12:02,145
15万到17万5的样子

1505
01:12:02,216 --> 01:12:04,656
你这么认为？

1506
01:12:04,726 --> 01:12:06,516
他们的确窃听了我们的电话

1507
01:12:06,586 --> 01:12:09,456
我的确这么认为，所以，呃，我们打算给你们...

1508
01:12:09,527 --> 01:12:10,957
我看我们该...

1509
01:12:11,027 --> 01:12:12,657
告诉开普勒法官这事

1510
01:12:12,727 --> 01:12:14,957
不，我反对

1511
01:12:15,038 --> 01:12:17,088
为什么？

1512
01:12:19,338 --> 01:12:20,928
我有个主意

1513
01:12:21,009 --> 01:12:22,099
呃...

1514
01:12:22,179 --> 01:12:25,439
有点疯狂的主意

1515
01:12:25,509 --> 01:12:26,979
妨碍司法公正？

1516
01:12:27,049 --> 01:12:28,279
我喜欢这样

1517
01:12:29,350 --> 01:12:30,940
我喜欢这样

1518
01:12:31,020 --> 01:12:33,080
我们最大的恶梦是谁呢？

1519
01:12:34,250 --> 01:12:35,950
法官给了我们陪审团的预选名单

1520
01:12:36,021 --> 01:12:37,961
一共92人

1521
01:12:38,031 --> 01:12:40,021
我们调查了他们的背景

1522
01:12:40,091 --> 01:12:42,291
确定了他们对我们有利还是不利

1523
01:12:42,362 --> 01:12:43,992
你打算让我怎么做？

1524
01:12:44,062 --> 01:12:45,462
好吧

1525
01:12:45,532 --> 01:12:46,502
什么？

1526
01:12:46,572 --> 01:12:48,032
任何与他们的直接接触

1527
01:12:48,103 --> 01:12:49,503
毫无疑问是在犯罪

1528
01:12:49,573 --> 01:12:50,973
我们在干嘛？

1529
01:12:51,043 --> 01:12:52,473
我们要这么做

1530
01:12:52,543 --> 01:12:53,803
照我说的做

1531
01:12:58,414 --> 01:12:59,474
喂

1532
01:12:59,554 --> 01:13:01,984
对，请给我找鲁迪.贝勒

1533
01:13:02,055 --> 01:13:03,985
我就是

1534
01:13:04,055 --> 01:13:06,925
我是比利.波特

1535
01:13:06,995 --> 01:13:08,825
你今天来我店里了

1536
01:13:08,896 --> 01:13:10,126
是的，波特先生

1537
01:13:10,196 --> 01:13:12,256
很感谢你给我回电话

1538
01:13:12,336 --> 01:13:14,996
你要干什么？

1539
01:13:15,067 --> 01:13:17,007
呃，和一个案子有关

1540
01:13:17,067 --> 01:13:19,167
你知道的，你被传召进入陪审团的那个

1541
01:13:19,237 --> 01:13:20,497
我是那个案子一方的律师

1542
01:13:20,577 --> 01:13:23,007
呃，这...这合法吗？

1543
01:13:23,078 --> 01:13:25,508
呃，当然，当然合法，波特先生

1544
01:13:25,578 --> 01:13:27,708
不过千万别告诉其他人

1545
01:13:27,778 --> 01:13:31,518
我代表那位女士，她孩子死了，是白血病

1546
01:13:31,589 --> 01:13:34,059
而这都是因为GB保险公司

1547
01:13:34,119 --> 01:13:35,089
不给他们钱

1548
01:13:35,160 --> 01:13:37,060
所以他们没法做手术

1549
01:13:37,130 --> 01:13:38,530
哦，天啊，这太恐怖了

1550
01:13:38,600 --> 01:13:41,860
我也有个婶婶得这个病

1551
01:13:41,931 --> 01:13:45,031
花了很多钱，我叔叔很痛苦...

1552
01:13:48,111 --> 01:13:49,701
我会尽力的

1553
01:13:49,772 --> 01:13:51,212
好的，先生

1554
01:13:51,272 --> 01:13:53,212
谢谢您，波特先生

1555
01:13:53,282 --> 01:13:54,972
再见

1556
01:13:55,052 --> 01:13:57,072
那个...那个婶婶是怎么回事？

1557
01:13:57,153 --> 01:13:58,553
那婶婶是哪冒出来的？

1558
01:13:58,623 --> 01:13:59,553
我怎么知道

1559
01:13:59,623 --> 01:14:01,053
你让我装得同情些嘛

1560
01:14:01,123 --> 01:14:03,093
我知道，可是你也不用...

1561
01:14:03,154 --> 01:14:04,094
不用说这么详细吧

1562
01:14:04,154 --> 01:14:05,054
拿着你的咖啡

1563
01:14:05,124 --> 01:14:06,564
我们回去

1564
01:14:07,824 --> 01:14:09,594
我只是想帮忙而已

1565
01:14:09,665 --> 01:14:10,565
快过来

1566
01:14:10,635 --> 01:14:12,195
我就在你后面

1567
01:14:14,165 --> 01:14:15,795
女士们先生们

1568
01:14:15,876 --> 01:14:18,366
我下面要问的问题

1569
01:14:18,446 --> 01:14:21,896
是今天最重要的问题

1570
01:14:21,976 --> 01:14:23,606
不过也很简单

1571
01:14:23,677 --> 01:14:25,837
用是或否回答就行

1572
01:14:25,917 --> 01:14:28,437
请听仔细

1573
01:14:28,517 --> 01:14:32,647
你们之中是否有人已经就此案被联系过了？

1574
01:14:34,628 --> 01:14:36,618
这是非常严重的事件

1575
01:14:36,699 --> 01:14:38,599
成功了

1576
01:14:38,669 --> 01:14:41,259
我们现在要知道...

1577
01:14:41,329 --> 01:14:43,929
让我换种方法提问吧

1578
01:14:46,110 --> 01:14:48,370
你们之中是否有人最近

1579
01:14:48,440 --> 01:14:51,540
和鲁迪.贝勒先生

1580
01:14:51,611 --> 01:14:54,241
或者他后面那位迭克.席夫列特先生谈过？

1581
01:14:54,311 --> 01:14:55,281
反对，法官大人！

1582
01:14:55,351 --> 01:14:56,581
这是侮辱！

1583
01:14:56,652 --> 01:14:58,242
你在做什么，德拉门德先生？

1584
01:14:58,322 --> 01:14:59,652
法官大人，

1585
01:14:59,722 --> 01:15:01,352
我们有证据表明

1586
01:15:01,422 --> 01:15:03,252
有人试图控制陪审团

1587
01:15:03,323 --> 01:15:04,553
他指的是我吧

1588
01:15:04,623 --> 01:15:07,923
我不理解你在干什么，德拉门德先生

1589
01:15:07,993 --> 01:15:08,933
我也不知道，法官大人

1590
01:15:09,003 --> 01:15:09,933
我也不知道

1591
01:15:10,004 --> 01:15:12,064
请上前

1592
01:15:14,834 --> 01:15:17,934
法官大人，有人试图控制陪审团

1593
01:15:18,005 --> 01:15:19,105
我要证据，立奥

1594
01:15:19,175 --> 01:15:21,975
我不能...

1595
01:15:22,045 --> 01:15:23,705
泄露秘密情报，法官大人

1596
01:15:23,786 --> 01:15:24,806
你疯了吗？

1597
01:15:24,886 --> 01:15:26,446
你看起来不正常

1598
01:15:26,516 --> 01:15:27,786
我能证明的

1599
01:15:27,856 --> 01:15:29,346
怎么证明？

1600
01:15:29,427 --> 01:15:31,447
控告我们这种罪名

1601
01:15:31,527 --> 01:15:32,827
控制陪审团

1602
01:15:32,897 --> 01:15:34,417
真可笑

1603
01:15:34,497 --> 01:15:36,467
请让我先问完陪审团

1604
01:15:36,528 --> 01:15:38,128
我相信您会明白事实的

1605
01:15:38,198 --> 01:15:40,728
你反对吗，贝勒先生？

1606
01:15:40,798 --> 01:15:41,928
不，我不反对

1607
01:15:42,008 --> 01:15:44,768
那好，继续

1608
01:15:44,839 --> 01:15:45,929
好

1609
01:15:48,349 --> 01:15:51,869
那到底怎么回事？

1610
01:15:51,950 --> 01:15:53,640
没什么，律师的小伎俩

1611
01:15:54,550 --> 01:15:55,750
波特先生

1612
01:15:55,820 --> 01:16:00,420
我想直接问你个问题

1613
01:16:00,491 --> 01:16:01,961
我希望你...

1614
01:16:02,031 --> 01:16:05,321
可以诚实的回答我

1615
01:16:05,402 --> 01:16:07,662
你只要诚实的问我

1616
01:16:07,732 --> 01:16:09,322
我就诚实的回答

1617
01:16:09,402 --> 01:16:11,562
很好

1618
01:16:11,643 --> 01:16:13,693
你有没有，波特先生，你几天前

1619
01:16:13,773 --> 01:16:15,933
难道没有和鲁迪.贝勒先生

1620
01:16:16,013 --> 01:16:17,373
通过电话吗？

1621
01:16:17,444 --> 01:16:18,574
废话，当然没有

1622
01:16:18,644 --> 01:16:20,974
我本以为你会给我个...

1623
01:16:21,044 --> 01:16:22,104
诚实的回答的

1624
01:16:22,184 --> 01:16:24,414
我就是实话实说

1625
01:16:24,485 --> 01:16:26,575
你肯定吗，波特先生？

1626
01:16:26,655 --> 01:16:28,315
我非常肯定

1627
01:16:28,395 --> 01:16:29,585
波特先生

1628
01:16:29,655 --> 01:16:31,025
在法庭上

1629
01:16:31,096 --> 01:16:32,786
在美国法庭上

1630
01:16:32,866 --> 01:16:35,126
在田纳西州公正的法庭上

1631
01:16:35,196 --> 01:16:37,226
我说你在说谎

1632
01:16:37,297 --> 01:16:38,637
你怎么敢说我说谎

1633
01:16:38,707 --> 01:16:39,667
德拉门德先生

1634
01:16:39,737 --> 01:16:42,327
你个婊子养的

1635
01:16:44,078 --> 01:16:45,598
嘿，把他拉出去！

1636
01:16:45,678 --> 01:16:46,868
肃静，保持法庭秩序！

1637
01:16:48,948 --> 01:16:50,418
贝立夫

1638
01:16:50,479 --> 01:16:52,349
把波特先生带出法庭

1639
01:16:52,419 --> 01:16:55,579
比利.波特先生，你被驱逐出陪审团了

1640
01:16:59,360 --> 01:17:00,450
法官大人

1641
01:17:00,530 --> 01:17:04,050
我要求解散整个陪审团

1642
01:17:04,131 --> 01:17:06,561
否决

1643
01:17:06,631 --> 01:17:07,961
陪审团已经被控制了

1644
01:17:14,712 --> 01:17:16,942
你...

1645
01:17:18,083 --> 01:17:20,883
你掉了你的鞋...

1646
01:17:23,824 --> 01:17:25,384
我们可以继续...

1647
01:17:25,454 --> 01:17:28,424
挑选陪审团了吗，德拉门德先生？

1648
01:17:28,494 --> 01:17:31,924
可以，法官大人

1649
01:17:31,995 --> 01:17:34,295
谢谢你

1650
01:17:57,328 --> 01:17:59,298
我们今天就提出诉讼

1651
01:17:59,359 --> 01:18:00,449
他会疯的

1652
01:18:00,529 --> 01:18:01,519
他已经疯了

1653
01:18:01,599 --> 01:18:04,529
他会干掉你的

1654
01:18:04,600 --> 01:18:05,830
我正盼着呢

1655
01:18:05,900 --> 01:18:08,630
这火鸡烤的不错

1656
01:18:08,710 --> 01:18:10,260
-哦
-喔

1657
01:18:10,340 --> 01:18:12,740
我，我得去法庭了

1658
01:18:12,811 --> 01:18:14,831
我已经迟到了

1659
01:18:14,911 --> 01:18:16,471
你的三明治怎么办？

1660
01:18:16,551 --> 01:18:19,171
我边走边吃吧

1661
01:18:19,252 --> 01:18:21,652
好吧

1662
01:18:21,722 --> 01:18:23,052
待会见

1663
01:18:23,122 --> 01:18:24,812
再见？

1664
01:18:33,604 --> 01:18:35,154
怎么样了？

1665
01:18:35,234 --> 01:18:37,134
克立佛去吃午饭时

1666
01:18:37,204 --> 01:18:38,724
我把文件给他了

1667
01:18:38,805 --> 01:18:40,365
他说他要打架

1668
01:18:40,445 --> 01:18:42,235
我说我有的是人打群架

1669
01:18:42,315 --> 01:18:43,335
他就软了

1670
01:18:43,415 --> 01:18:46,505
不过，老兄，你最好小心点

1671
01:18:48,556 --> 01:18:49,916
好的

1672
01:18:49,986 --> 01:18:53,186
谢谢你，布奇，非常感谢

1673
01:18:53,257 --> 01:18:54,557
你庄严宣誓

1674
01:18:54,627 --> 01:18:57,027
你下面要做的证供

1675
01:18:57,097 --> 01:18:58,117
是事实，并且是全部事实

1676
01:18:58,198 --> 01:19:00,098
没有丝毫谎言吗？

1677
01:19:00,168 --> 01:19:01,128
是的，我宣誓

1678
01:19:01,198 --> 01:19:03,968
你可以作证了

1679
01:19:10,339 --> 01:19:13,309
请告诉法庭你的名字，以便记录

1680
01:19:13,380 --> 01:19:17,080
玛嘉丽.布莱克夫人

1681
01:19:17,150 --> 01:19:18,740
好，布莱克夫人

1682
01:19:18,821 --> 01:19:21,051
你是丹尼.雷.布莱克的母亲

1683
01:19:21,121 --> 01:19:23,381
你儿子是否最近因白血病而死去

1684
01:19:23,461 --> 01:19:24,621
而这些都是因为，被告，GB保险公司...

1685
01:19:24,691 --> 01:19:26,021
反对

1686
01:19:26,092 --> 01:19:28,032
诱导证人

1687
01:19:28,092 --> 01:19:31,292
反对有效

1688
01:19:33,833 --> 01:19:34,933
你的儿子

1689
01:19:35,003 --> 01:19:37,063
丹尼雷，需要做手术

1690
01:19:37,143 --> 01:19:38,073
反对

1691
01:19:38,143 --> 01:19:39,233
诱导证人

1692
01:19:39,314 --> 01:19:40,934
反对有效

1693
01:19:43,444 --> 01:19:46,174
布莱克夫人，你买这个医疗保险

1694
01:19:46,255 --> 01:19:48,415
是不是想要为你的儿子提供医疗保障？

1695
01:19:48,485 --> 01:19:49,415
反对

1696
01:19:49,485 --> 01:19:50,855
对不起，法官大人

1697
01:19:50,925 --> 01:19:52,325
诱导证人

1698
01:19:53,326 --> 01:19:54,486
贝勒先生

1699
01:19:54,556 --> 01:19:56,286
你干嘛不出示一下保单

1700
01:19:56,366 --> 01:19:58,026
然后问问她干嘛要买这个？

1701
01:19:58,096 --> 01:19:59,356
好的

1702
01:20:01,337 --> 01:20:03,857
贝勒先生

1703
01:20:03,937 --> 01:20:06,767
你必须先得到许可才能上前盘问证人

1704
01:20:06,838 --> 01:20:10,938
对不起，法官大人

1705
01:20:11,008 --> 01:20:12,308
请求上前盘问证人

1706
01:20:12,379 --> 01:20:13,579
要求批准

1707
01:20:17,019 --> 01:20:19,219
放轻松点

1708
01:20:19,290 --> 01:20:21,920
放轻松点

1709
01:20:21,990 --> 01:20:25,190
“GB保险公司，1996年7月7日，

1710
01:20:25,260 --> 01:20:33,791
回复：7849909886号保单

1711
01:20:33,872 --> 01:20:36,132
亲爱的布莱克夫人

1712
01:20:36,212 --> 01:20:38,372
本公司已经7次

1713
01:20:38,442 --> 01:20:42,242
以书面方式拒绝了你的赔付要求

1714
01:20:42,313 --> 01:20:47,883
这是第八次拒绝，也是最后一次

1715
01:20:47,954 --> 01:20:52,914
你实在是蠢，蠢，蠢！

1716
01:20:54,125 --> 01:20:56,855
诚挚的，埃佛莱特.拉佛金，

1717
01:20:56,935 --> 01:20:59,295
赔付副总裁”

1718
01:21:02,306 --> 01:21:03,396
请再读一遍

1719
01:21:03,476 --> 01:21:04,596
反对

1720
01:21:04,676 --> 01:21:05,636
无意义的重复，法官大人

1721
01:21:05,706 --> 01:21:06,766
反对有效

1722
01:21:09,417 --> 01:21:11,107
我问完了

1723
01:21:11,177 --> 01:21:12,737
德拉门德先生？

1724
01:21:12,818 --> 01:21:14,718
法官大人

1725
01:21:14,788 --> 01:21:17,878
请收起那个投影

1726
01:21:17,958 --> 01:21:20,388
现在，布莱克夫人

1727
01:21:23,029 --> 01:21:24,519
把它关掉

1728
01:21:28,700 --> 01:21:33,300
布莱克夫人，你为什么要控告GB保险公司，
索赔1000万美元呢？

1729
01:21:33,371 --> 01:21:35,431
就这么点吗？

1730
01:21:36,941 --> 01:21:38,501
对不起？

1731
01:21:38,581 --> 01:21:41,101
我认为我们该要得更多呢

1732
01:21:41,182 --> 01:21:42,912
是吗？

1733
01:21:42,982 --> 01:21:45,252
是的，先生，你的客户有十亿美元

1734
01:21:45,322 --> 01:21:47,182
而且你的客户杀了我的儿子

1735
01:21:47,253 --> 01:21:50,853
我打算索赔更多

1736
01:21:51,993 --> 01:21:54,623
你拿这笔钱打算怎么办？

1737
01:21:54,694 --> 01:21:56,594
如果陪审团判给你1000万

1738
01:21:56,664 --> 01:21:58,104
你打算怎么花？

1739
01:21:58,164 --> 01:21:59,104
我会把它们捐给

1740
01:21:59,174 --> 01:22:02,034
美国白血病研究学会

1741
01:22:02,105 --> 01:22:03,735
全部捐掉

1742
01:22:04,975 --> 01:22:08,505
我不要你的臭钱

1743
01:22:08,576 --> 01:22:12,606
请注意你是宣誓过的，布莱克夫人

1744
01:22:12,686 --> 01:22:14,846
请求上前盘问证人，法官大人？

1745
01:22:14,917 --> 01:22:16,187
要求批准

1746
01:22:16,257 --> 01:22:17,227
布莱克夫人

1747
01:22:17,287 --> 01:22:19,227
我希望你可以读一下保单

1748
01:22:19,297 --> 01:22:23,557
第16页，第K项，第14段，第E款

1749
01:22:23,628 --> 01:22:24,928
保险公司

1750
01:22:24,998 --> 01:22:26,558
用浅显的英语说

1751
01:22:26,638 --> 01:22:29,938
它不负责赔付实验性治疗

1752
01:22:30,009 --> 01:22:32,739
你的起诉书中说

1753
01:22:32,809 --> 01:22:35,069
如果你的儿子可以接受骨髓移植

1754
01:22:35,140 --> 01:22:36,480
他就能够得救

1755
01:22:36,550 --> 01:22:38,540
布莱克夫人

1756
01:22:38,610 --> 01:22:40,010
难道美国不是

1757
01:22:40,080 --> 01:22:42,180
一年只有7000次

1758
01:22:42,251 --> 01:22:43,581
这种骨髓移植吗？

1759
01:22:43,651 --> 01:22:45,251
而田纳西州只有不到200次

1760
01:22:45,321 --> 01:22:46,481
反对，法官大人

1761
01:22:46,561 --> 01:22:47,751
他诱导证人

1762
01:22:47,822 --> 01:22:50,262
这是交叉盘问，诱导是允许的

1763
01:22:50,332 --> 01:22:52,122
就诱导而言

1764
01:22:52,202 --> 01:22:54,102
反对无效

1765
01:22:54,163 --> 01:22:58,603
所以这不是保单赔付范围

1766
01:22:58,673 --> 01:22:59,663
还有，布莱克夫人

1767
01:22:59,743 --> 01:23:02,473
谁是...

1768
01:23:02,544 --> 01:23:05,444
谁第一个诊断你儿子病情的

1769
01:23:07,085 --> 01:23:09,015
最开始

1770
01:23:09,085 --> 01:23:12,545
是贝奇医生

1771
01:23:12,625 --> 01:23:14,985
你的家庭医生？

1772
01:23:15,056 --> 01:23:16,686
是的，先生

1773
01:23:16,756 --> 01:23:18,156
他是个好医生吗？

1774
01:23:18,226 --> 01:23:19,926
他是个很好的医生

1775
01:23:19,996 --> 01:23:22,486
那么，布莱克夫人

1776
01:23:22,567 --> 01:23:23,997
这个诚恳而能力卓著的医生

1777
01:23:24,067 --> 01:23:26,967
是不是反覆告诉过你，

1778
01:23:27,037 --> 01:23:29,057
因为你儿子所患的特殊白血病类型

1779
01:23:29,138 --> 01:23:31,908
骨髓移植对你儿子的病情不利？

1780
01:23:33,478 --> 01:23:35,278
不，没有

1781
01:23:35,349 --> 01:23:36,649
他没有说过

1782
01:23:36,719 --> 01:23:37,679
他没这么说过

1783
01:23:37,749 --> 01:23:40,309
他不是那么对我说的

1784
01:23:41,820 --> 01:23:42,980
请求上前盘问证人，法官大人？

1785
01:23:43,060 --> 01:23:44,450
要求批准

1786
01:23:44,520 --> 01:23:50,220
布莱克夫人，这是不是贝奇医生的信签？

1787
01:23:50,301 --> 01:23:53,291
他是不是在底下签字了？

1788
01:23:53,371 --> 01:23:55,131
他不能这么做

1789
01:23:55,202 --> 01:23:56,602
为什么？

1790
01:23:56,672 --> 01:23:57,662
因为他不能这么

1791
01:23:57,742 --> 01:23:59,972
提出证据，更何况那只是道听途说

1792
01:24:00,042 --> 01:24:01,872
反对，法官大人

1793
01:24:01,943 --> 01:24:05,843
布莱克家庭医生给德拉门德先生的

1794
01:24:05,913 --> 01:24:07,973
信是不可接受的

1795
01:24:08,054 --> 01:24:09,954
很正确，法官大人

1796
01:24:10,024 --> 01:24:11,854
我没有要求此信

1797
01:24:11,924 --> 01:24:13,324
被作为呈堂证供

1798
01:24:13,394 --> 01:24:16,014
我只是想知道，证人是否...

1799
01:24:16,095 --> 01:24:17,155
可以阅读此信

1800
01:24:17,225 --> 01:24:18,955
这是根据田纳西州

1801
01:24:19,035 --> 01:24:20,295
取证法第612条

1802
01:24:20,365 --> 01:24:22,165
是为了让她恢复回忆

1803
01:24:22,236 --> 01:24:24,666
贝勒先生，你觉得呢？

1804
01:24:24,736 --> 01:24:26,296
我不知道，法官大人

1805
01:24:26,376 --> 01:24:27,496
但我反对这样

1806
01:24:27,506 --> 01:24:31,266
而且在预审举证时

1807
01:24:31,337 --> 01:24:32,737
这封信并未被提出

1808
01:24:32,807 --> 01:24:35,277
对此你有何解释，德拉门德先生？

1809
01:24:35,348 --> 01:24:37,078
我不知道我们会需要这封信

1810
01:24:37,148 --> 01:24:38,848
我本来以为这位女士

1811
01:24:38,918 --> 01:24:40,348
会真实的告诉我们她的医生对她说了什么

1812
01:24:42,559 --> 01:24:44,649
还有什么事吗，贝勒先生？

1813
01:24:46,789 --> 01:24:48,259
没了

1814
01:24:50,400 --> 01:24:51,560
德拉门德先生

1815
01:24:51,630 --> 01:24:54,290
我不是反对你这么做

1816
01:24:54,370 --> 01:24:55,990
但是最好别扯太远了

1817
01:24:56,071 --> 01:24:58,631
没问题，法官大人。现在，布莱克夫人

1818
01:24:58,701 --> 01:25:01,261
这封信有没有让你回忆起

1819
01:25:01,341 --> 01:25:03,571
丹尼雷的那种白血病类型

1820
01:25:03,642 --> 01:25:04,942
到底能不能通过

1821
01:25:05,012 --> 01:25:06,172
骨髓移植治疗呢？

1822
01:25:06,252 --> 01:25:08,542
可是，他不是专家

1823
01:25:08,613 --> 01:25:09,553
但是他是执照医师

1824
01:25:09,623 --> 01:25:12,883
很有经验，年富力强

1825
01:25:12,953 --> 01:25:15,683
而且反覆告诉过你

1826
01:25:15,764 --> 01:25:18,954
你理智上无法接受的事情：

1827
01:25:19,034 --> 01:25:21,864
不管怎么治疗

1828
01:25:21,935 --> 01:25:24,735
你的儿子仍然会因为白血病而死

1829
01:25:24,805 --> 01:25:26,165
是不是？

1830
01:25:26,235 --> 01:25:28,795
但是他不是这方面的专家

1831
01:25:28,876 --> 01:25:30,276
我不相信他

1832
01:25:30,346 --> 01:25:33,176
不管你相不相信他，夫人

1833
01:25:33,246 --> 01:25:35,106
我不相信你！

1834
01:25:35,177 --> 01:25:36,847
你还记得就在几分钟前

1835
01:25:36,917 --> 01:25:38,747
你经过宣誓后对陪审团说的话吗？

1836
01:25:38,817 --> 01:25:41,587
你说贝奇医生从来没有说过

1837
01:25:41,657 --> 01:25:44,087
你儿子的白血病

1838
01:25:44,158 --> 01:25:45,558
无法通过骨髓移植治疗。

1839
01:25:45,628 --> 01:25:47,688
我相信你的原话是：

1840
01:25:47,758 --> 01:25:49,918
“他没这么说过

1841
01:25:49,999 --> 01:25:51,129
他不是那么...

1842
01:25:51,199 --> 01:25:52,959
对我说的”

1843
01:25:53,039 --> 01:25:55,939
他不是专家

1844
01:25:56,010 --> 01:25:58,300
我只是希望丹尼雷

1845
01:25:58,370 --> 01:26:03,540
可以接受最好的治疗

1846
01:26:03,611 --> 01:26:06,141
你也会这么做的

1847
01:26:06,221 --> 01:26:08,271
当然，夫人

1848
01:26:08,351 --> 01:26:11,051
当然

1849
01:26:11,122 --> 01:26:12,682
我问完了

1850
01:26:15,062 --> 01:26:17,252
你可以退席了，布莱克夫人

1851
01:26:24,874 --> 01:26:27,304
我是不是做的不好？

1852
01:26:27,374 --> 01:26:28,964
不，不，你做的很好

1853
01:26:29,045 --> 01:26:30,135
很好

1854
01:26:30,215 --> 01:26:32,275
没关系，陪审团会

1855
01:26:32,345 --> 01:26:33,905
明确知道他想干什么的

1856
01:26:33,985 --> 01:26:35,385
陪审团会知道的

1857
01:26:35,445 --> 01:26:36,415
我想抽烟

1858
01:26:36,486 --> 01:26:37,386
我知道你需要

1859
01:26:37,456 --> 01:26:38,286
我们过会再抽

1860
01:26:41,986 --> 01:26:44,926
我希望他没有换过锁

1861
01:26:44,997 --> 01:26:46,557
你害怕吗？

1862
01:26:46,627 --> 01:26:48,567
是的

1863
01:26:50,768 --> 01:26:52,468
我们干吧

1864
01:27:19,302 --> 01:27:21,032
和猪圈一样

1865
01:27:21,972 --> 01:27:23,162
对不起

1866
01:27:23,243 --> 01:27:25,903
快点，凯莉，快点

1867
01:27:28,643 --> 01:27:31,043
我的东西基本都在衣橱里

1868
01:27:39,555 --> 01:27:41,995
这的东西你是拿不光的，凯莉

1869
01:27:57,448 --> 01:27:59,508
哦，不

1870
01:28:00,448 --> 01:28:02,938
嘿，你好啊！

1871
01:28:03,018 --> 01:28:04,538
我回家啦

1872
01:28:05,449 --> 01:28:07,349
哇，这是什么人啊？

1873
01:28:07,419 --> 01:28:09,979
看看什么人来了？

1874
01:28:10,060 --> 01:28:11,750
你们俩在这里干嘛，啊？

1875
01:28:11,830 --> 01:28:13,520
嘿，放轻松点，好吗？

1876
01:28:13,600 --> 01:28:14,530
过来，告诉我

1877
01:28:14,600 --> 01:28:15,590
我是你的丈夫，还记得吗，啊？

1878
01:28:15,660 --> 01:28:16,600
老兄，放轻松点

1879
01:28:16,670 --> 01:28:17,560
闭嘴！

1880
01:28:17,631 --> 01:28:19,331
嘿，听着，没事的

1881
01:28:19,401 --> 01:28:20,331
我才不听你废话呢！

1882
01:28:20,401 --> 01:28:22,201
放松点

1883
01:28:22,271 --> 01:28:23,201
啊！

1884
01:28:23,271 --> 01:28:24,171
你让我心痛啊，宝贝

1885
01:28:24,242 --> 01:28:25,402
你让我很难受

1886
01:28:26,542 --> 01:28:28,912
你干嘛要这么对我？!

1887
01:28:28,942 --> 01:28:29,542
不，克立佛!

1888
01:28:33,953 --> 01:28:35,683
快跑！

1889
01:29:08,475 --> 01:29:10,145
鲁迪，鲁迪！

1890
01:29:10,215 --> 01:29:11,205
克立佛！

1891
01:29:11,286 --> 01:29:12,346
你这个笨蛋，你在干嘛？!

1892
01:29:12,416 --> 01:29:13,386
你看到了吧？

1893
01:29:13,456 --> 01:29:14,976
凯莉！

1894
01:29:15,056 --> 01:29:16,986
这就是你想要的？

1895
01:29:17,056 --> 01:29:18,186
不是我的错哦！

1896
01:29:18,257 --> 01:29:20,227
我是爱你的！我爱你！

1897
01:29:32,979 --> 01:29:34,529
鲁迪！

1898
01:29:40,220 --> 01:29:41,770
别打了，鲁迪

1899
01:29:41,850 --> 01:29:43,340
别打

1900
01:29:43,420 --> 01:29:45,050
哦

1901
01:29:45,121 --> 01:29:47,211
给我那根棒子，快走

1902
01:29:47,291 --> 01:29:49,021
什么？

1903
01:29:50,861 --> 01:29:53,191
给我那根棒子，快走

1904
01:29:53,262 --> 01:29:55,862
你今天晚上没有来过

1905
01:29:55,932 --> 01:29:57,662
给我那根棒子

1906
01:30:12,225 --> 01:30:13,185
快走，鲁迪

1907
01:30:13,255 --> 01:30:15,225
你今天晚上没有来过

1908
01:30:52,331 --> 01:30:54,231
他最终杀了她

1909
01:30:54,301 --> 01:30:56,231
不，是他，她把他杀了

1910
01:30:56,301 --> 01:30:57,331
你肯定？

1911
01:30:57,401 --> 01:30:58,371
我刚刚看过她的。

1912
01:30:58,442 --> 01:30:59,602
怎么发生的？

1913
01:30:59,672 --> 01:31:00,662
我不知道

1914
01:31:00,742 --> 01:31:02,442
有人说当你谋杀一个人时,

1915
01:31:02,512 --> 01:31:03,812
你起码犯了25项错误

1916
01:31:03,882 --> 01:31:07,002
过后，你要是能记住其中5项就不错了

1917
01:31:07,083 --> 01:31:08,143
这是自卫行为

1918
01:31:08,213 --> 01:31:10,243
但是我没法忘记他死了

1919
01:31:10,323 --> 01:31:12,753
这个错误充斥了我的大脑

1920
01:31:12,824 --> 01:31:14,484
我没法思考了

1921
01:31:16,294 --> 01:31:17,484
但是凯莉知道

1922
01:31:17,564 --> 01:31:19,124
她知道该做什么

1923
01:31:19,195 --> 01:31:20,725
她知道这是好机会

1924
01:31:20,795 --> 01:31:23,595
当那些事情

1925
01:31:23,665 --> 01:31:25,435
一起发生的时候

1926
01:31:25,506 --> 01:31:29,236
她只想到我的安全

1927
01:31:29,306 --> 01:31:32,146
而我把她一个人扔在那里

1928
01:31:33,647 --> 01:31:35,447
他妈的！

1929
01:31:35,517 --> 01:31:36,707
他妈的，凯莉！

1930
01:31:36,787 --> 01:31:38,447
这到底发生什么了？

1931
01:31:38,517 --> 01:31:40,457
你到底做什么了？

1932
01:31:47,929 --> 01:31:49,559
你杀了我的儿子！

1933
01:31:49,629 --> 01:31:50,719
你他妈的！

1934
01:32:25,704 --> 01:32:28,104
对不起，我是她的律师

1935
01:32:28,175 --> 01:32:31,205
我要求在她被问讯时在场

1936
01:32:31,275 --> 01:32:32,675
这是你的律师？

1937
01:32:32,746 --> 01:32:34,306
是的，长官

1938
01:32:36,216 --> 01:32:39,156
我希望能够保释她，由我监管她

1939
01:32:39,226 --> 01:32:40,386
我不能这样做

1940
01:32:40,457 --> 01:32:42,357
我不知道你是哪种律师

1941
01:32:42,427 --> 01:32:43,917
不过这是人命官司

1942
01:32:43,997 --> 01:32:46,857
只有法官才能作出是否保释的决定

1943
01:32:46,928 --> 01:32:48,458
我会进监狱吗？

1944
01:32:48,528 --> 01:32:51,058
你们能给她安排个单间吗？

1945
01:32:51,138 --> 01:32:53,368
我说，小子，我又不是监狱长

1946
01:32:53,439 --> 01:32:55,409
你干嘛不做点更有效的事情呢

1947
01:32:55,469 --> 01:32:57,599
和监狱官谈谈？

1948
01:32:57,679 --> 01:32:59,909
他们可喜欢律师了，是不？

1949
01:32:59,980 --> 01:33:02,410
我说...

1950
01:33:02,480 --> 01:33:04,420
要是你的律师有用的话

1951
01:33:04,480 --> 01:33:05,640
只要你能交保释金

1952
01:33:05,720 --> 01:33:07,910
明天什么时候你就能出狱了

1953
01:33:10,021 --> 01:33:11,961
好吧

1954
01:33:16,302 --> 01:33:17,492
你们有5分钟时间谈话

1955
01:33:17,562 --> 01:33:18,762
谢谢你

1956
01:33:25,043 --> 01:33:27,033
他们从窗口那里看着我们

1957
01:33:27,114 --> 01:33:29,804
而且这个房间可能安了窃听器

1958
01:33:29,884 --> 01:33:32,904
所以你要小心你说的话

1959
01:33:32,984 --> 01:33:34,814
"一般杀人罪"什么意思？

1960
01:33:34,885 --> 01:33:37,755
就是说非故意的杀人

1961
01:33:37,825 --> 01:33:39,685
我还能活多久？

1962
01:33:39,755 --> 01:33:41,285
不，不，你得先被定罪才可能判刑

1963
01:33:41,356 --> 01:33:43,886
而我不会让你定罪的

1964
01:33:43,966 --> 01:33:46,726
那绝不可能发生

1965
01:33:50,267 --> 01:33:53,827
夫人，请把双手交叉在身后放好

1966
01:33:58,708 --> 01:34:00,678
这边走，夫人

1967
01:34:07,620 --> 01:34:09,150
拉佛金先生

1968
01:34:09,230 --> 01:34:11,090
你是GB保险公司的

1969
01:34:11,160 --> 01:34:12,490
赔付副总裁，是不是？

1970
01:34:12,560 --> 01:34:13,590
是的

1971
01:34:13,661 --> 01:34:14,491
请求上前盘问证人，法官大人？

1972
01:34:14,561 --> 01:34:15,861
要求批准

1973
01:34:18,101 --> 01:34:20,501
你认得这个吗？

1974
01:34:20,572 --> 01:34:21,632
继续

1975
01:34:23,812 --> 01:34:25,902
读过陪审团听

1976
01:34:28,683 --> 01:34:29,773
“亲爱的布莱克夫人,

1977
01:34:29,853 --> 01:34:31,183
本公司已经7次

1978
01:34:31,253 --> 01:34:33,113
以书面方式拒绝了你的赔付要求

1979
01:34:33,183 --> 01:34:34,813
这是第八次拒绝，也是最后一次

1980
01:34:34,884 --> 01:34:36,944
你实在是蠢，蠢，蠢！

1981
01:34:37,024 --> 01:34:38,684
诚挚的，埃佛莱特.拉佛金，

1982
01:34:38,754 --> 01:34:40,594
赔付副总裁”

1983
01:34:40,665 --> 01:34:42,655
这是你写的？

1984
01:34:42,735 --> 01:34:43,925
是的

1985
01:34:46,365 --> 01:34:47,995
你打算作何解释？

1986
01:34:53,446 --> 01:34:56,066
那是我个人遇到了很多困难

1987
01:34:56,147 --> 01:34:58,667
我有很多压力

1988
01:34:58,747 --> 01:35:02,237
我们已经7次拒绝这个赔付了

1989
01:35:02,318 --> 01:35:04,918
所以我想要用点比较强的词

1990
01:35:04,988 --> 01:35:06,578
我有点过分了

1991
01:35:06,658 --> 01:35:08,178
我后悔写了这样的信

1992
01:35:08,259 --> 01:35:10,689
我愿意道歉

1993
01:35:10,759 --> 01:35:11,989
你不觉得现在道歉

1994
01:35:12,059 --> 01:35:13,259
有点太晚了吗？

1995
01:35:15,070 --> 01:35:16,800
也许

1996
01:35:16,870 --> 01:35:17,840
也许？

1997
01:35:17,900 --> 01:35:19,870
那个孩子已经死了，是不是？

1998
01:35:23,311 --> 01:35:25,901
是的

1999
01:35:25,981 --> 01:35:27,911
那么，拉佛金先生

2000
01:35:27,982 --> 01:35:30,282
谁是杰琪.莱曼切克？

2001
01:35:30,352 --> 01:35:36,252
杰琪.莱曼切克是前赔付审核员

2002
01:35:36,323 --> 01:35:38,153
她在你的部门工作？

2003
01:35:38,223 --> 01:35:39,163
是的

2004
01:35:39,233 --> 01:35:41,323
她什么时候离开

2005
01:35:41,394 --> 01:35:42,364
GB保险公司的？

2006
01:35:42,434 --> 01:35:45,194
我记不住准确日期

2007
01:35:45,264 --> 01:35:47,204
是不是10月30日？

2008
01:35:47,274 --> 01:35:48,704
好像差不多

2009
01:35:48,775 --> 01:35:51,035
那是不是她本该就此案作证的

2010
01:35:51,105 --> 01:35:52,975
两天前？

2011
01:35:53,045 --> 01:35:54,605
我真的记不清

2012
01:35:54,676 --> 01:35:57,376
我希望根据第612条法令

2013
01:35:57,446 --> 01:36:00,116
让证人恢复一下回忆

2014
01:36:02,787 --> 01:36:04,157
是10月30日

2015
01:36:04,227 --> 01:36:05,347
很明显嘛

2016
01:36:05,427 --> 01:36:06,587
而这正是她本该就此案作证的

2017
01:36:06,657 --> 01:36:07,687
两天前？

2018
01:36:07,757 --> 01:36:09,387
我想是吧

2019
01:36:09,468 --> 01:36:10,858
而且正是她

2020
01:36:10,928 --> 01:36:12,418
处理丹尼雷的赔付申请的，

2021
01:36:12,498 --> 01:36:13,398
是不是？

2022
01:36:13,468 --> 01:36:15,268
是的

2023
01:36:15,339 --> 01:36:16,629
你把她辞掉了？

2024
01:36:18,669 --> 01:36:20,799
当然不是

2025
01:36:20,879 --> 01:36:23,209
那么，你怎么把她弄跑的？

2026
01:36:23,280 --> 01:36:24,370
她自己要求辞职的

2027
01:36:24,450 --> 01:36:26,920
就在你给我的信中就能看到

2028
01:36:26,980 --> 01:36:28,510
哦，那么她干嘛要辞职？

2029
01:36:30,091 --> 01:36:32,751
“我在此因个人原因辞职”

2030
01:36:32,821 --> 01:36:34,261
所以是她自愿

2031
01:36:34,321 --> 01:36:35,951
离开的？

2032
01:36:36,032 --> 01:36:37,122
她就是那么写的

2033
01:36:37,192 --> 01:36:39,252
没有要问的了

2034
01:36:39,332 --> 01:36:41,322
你可以退席了，先生

2035
01:36:48,203 --> 01:36:49,603
你好

2036
01:36:49,674 --> 01:36:54,514
嗨，我是杰琪.莱曼切克的兄弟，我叫詹姆斯

2037
01:36:54,584 --> 01:36:56,014
我能看看她吗？

2038
01:36:56,085 --> 01:36:57,205
詹姆斯.莱曼切克？

2039
01:36:57,285 --> 01:36:58,515
是的

2040
01:36:58,585 --> 01:37:00,285
稍等

2041
01:37:18,938 --> 01:37:20,538
全体起立

2042
01:37:22,949 --> 01:37:26,039
请让我解释一下，莱曼切克小姐

2043
01:37:26,119 --> 01:37:28,049
其实我不是你的兄弟

2044
01:37:32,160 --> 01:37:34,750
鲁迪

2045
01:37:34,830 --> 01:37:37,420
来见见杰琪.莱曼切克

2046
01:37:45,172 --> 01:37:48,002
这是杰琪.莱曼切克。

2047
01:37:48,072 --> 01:37:51,442
这是卡尔，他一定要跟着她

2048
01:37:53,113 --> 01:37:54,743
这是我的合伙人

2049
01:37:54,813 --> 01:37:56,583
鲁迪.S.贝勒

2050
01:37:56,654 --> 01:37:58,884
请把你告诉过我的再和他说一遍

2051
01:37:58,954 --> 01:38:00,984
莱曼切克小姐

2052
01:38:01,054 --> 01:38:02,824
很高兴见到你

2053
01:38:02,895 --> 01:38:04,955
我可以坐下吗

2054
01:38:05,025 --> 01:38:06,225
当然

2055
01:38:06,295 --> 01:38:08,785
很好，莱曼切克小姐

2056
01:38:08,865 --> 01:38:11,465
我想和你谈谈关于布莱克的赔付申请

2057
01:38:11,536 --> 01:38:13,166
那是指定由你处理的？

2058
01:38:13,236 --> 01:38:15,296
是的

2059
01:38:15,376 --> 01:38:17,606
最开始布莱克夫人的申请

2060
01:38:17,677 --> 01:38:18,697
是我处理的

2061
01:38:18,777 --> 01:38:20,337
根据公司当时的规定，

2062
01:38:20,417 --> 01:38:23,177
我给了她一封拒绝信

2063
01:38:23,248 --> 01:38:24,238
为什么？

2064
01:38:24,318 --> 01:38:26,248
为什么？

2065
01:38:27,988 --> 01:38:30,248
因为所有的赔付请求一开始都得被拒绝

2066
01:38:33,929 --> 01:38:35,419
所有的？

2067
01:38:35,499 --> 01:38:37,989
所有的

2068
01:38:43,471 --> 01:38:46,501
听着，他们是这么干的

2069
01:38:46,581 --> 01:38:48,771
他们...

2070
01:38:48,841 --> 01:38:50,541
一家家的在

2071
01:38:50,612 --> 01:38:52,052
穷人区卖保单

2072
01:38:52,112 --> 01:38:53,712
每星期有人

2073
01:38:53,782 --> 01:38:56,312
来收现金

2074
01:38:56,393 --> 01:39:01,293
如果有赔付申请，会有专门的审核员处理

2075
01:39:01,363 --> 01:39:05,623
不过那只不过是低级的打字工作

2076
01:39:05,694 --> 01:39:10,364
但不管如何，审核员看一遍申请后

2077
01:39:10,435 --> 01:39:13,205
立刻给投保人发信...

2078
01:39:13,275 --> 01:39:14,365
拒绝赔付

2079
01:39:14,445 --> 01:39:16,175
赔付申请的审核员

2080
01:39:16,245 --> 01:39:19,175
随即把文件发给承销人

2081
01:39:19,246 --> 01:39:21,686
承销人接着会给审核部门发一份备忘录

2082
01:39:21,746 --> 01:39:24,746
告诉他们，“如果没有我们的同意，不要赔付”

2083
01:39:24,817 --> 01:39:27,347
对了，你要记住

2084
01:39:27,427 --> 01:39:28,857
虽然这些人都是...

2085
01:39:28,927 --> 01:39:29,887
为这个大公司服务的

2086
01:39:29,958 --> 01:39:31,548
即使他们在同一间大楼内工作

2087
01:39:31,628 --> 01:39:33,118
他们却互不相识

2088
01:39:33,198 --> 01:39:35,718
也不知道彼此的部门打算怎么办

2089
01:39:35,798 --> 01:39:36,888
所有这些都是精心安排的

2090
01:39:36,969 --> 01:39:39,159
部门之间互相敌视，互不关心

2091
01:39:39,239 --> 01:39:41,329
同时，

2092
01:39:41,409 --> 01:39:44,269
保险的客户们...

2093
01:39:45,640 --> 01:39:48,610
不停的接到这种拒付信

2094
01:39:48,680 --> 01:39:52,670
有些来自审核部门，有些来自承销商

2095
01:39:52,751 --> 01:39:55,691
大多数人就此放弃了

2096
01:39:57,362 --> 01:40:01,052
而这些，当然...

2097
01:40:01,132 --> 01:40:02,862
也是经过精心设计的

2098
01:40:15,314 --> 01:40:16,744
你的下一个证人，贝勒先生

2099
01:40:18,115 --> 01:40:20,205
原告方传召杰琪.莱曼切克

2100
01:40:20,285 --> 01:40:22,115
他说什么？

2101
01:40:24,056 --> 01:40:25,116
哦

2102
01:40:26,826 --> 01:40:28,916
反对，法官大人！

2103
01:40:28,996 --> 01:40:29,956
我可以过来和您谈谈吗？

2104
01:40:30,026 --> 01:40:31,616
可以

2105
01:40:31,697 --> 01:40:33,597
这是突然袭击，法官大人

2106
01:40:33,667 --> 01:40:35,397
为什么？她被列为可能传召的证人了

2107
01:40:35,467 --> 01:40:36,497
但我们应该提前得到通知的

2108
01:40:36,567 --> 01:40:37,697
你什么时候找到她的？

2109
01:40:37,768 --> 01:40:38,928
我从来不知道她失踪了啊

2110
01:40:39,008 --> 01:40:40,498
你得回答这个问题，贝勒先生

2111
01:40:42,138 --> 01:40:43,128
这是我第一个案子

2112
01:40:43,208 --> 01:40:44,538
那也不行

2113
01:40:44,609 --> 01:40:46,879
这是个事关公平与否的问题

2114
01:40:46,949 --> 01:40:49,179
我们应该事前知道你要传召的证人

2115
01:40:49,249 --> 01:40:50,549
我同意

2116
01:40:50,619 --> 01:40:52,019
你难道是说她不能出庭作证？

2117
01:40:52,090 --> 01:40:54,780
在预审名单中有她的名字，法官大人

2118
01:40:54,860 --> 01:40:57,160
根据第26条法律第6款

2119
01:40:57,230 --> 01:41:00,720
我们有权传召她为证人

2120
01:41:02,401 --> 01:41:03,921
反对无效

2121
01:41:09,812 --> 01:41:11,362
去背一遍吧

2122
01:41:23,524 --> 01:41:25,544
请告诉法庭你的姓名，以便记录

2123
01:41:25,625 --> 01:41:27,325
杰琪.莱曼切克

2124
01:41:28,265 --> 01:41:29,225
莱曼切克小姐

2125
01:41:29,295 --> 01:41:31,195
你为GB保险公司工作了多久？

2126
01:41:31,266 --> 01:41:32,996
6年

2127
01:41:33,066 --> 01:41:34,796
你什么时候不再干这个工作的？

2128
01:41:34,866 --> 01:41:36,896
10月30日

2129
01:41:36,976 --> 01:41:38,166
你为什么不干了？

2130
01:41:38,237 --> 01:41:39,727
我被辞退了

2131
01:41:39,807 --> 01:41:42,367
你是说你不是辞职的？

2132
01:41:42,447 --> 01:41:44,247
不，我是被辞退的

2133
01:41:46,018 --> 01:41:47,348
我可以上前盘问证人吗，法官大人？

2134
01:41:47,418 --> 01:41:48,348
要求批准

2135
01:41:48,418 --> 01:41:49,788
我有些疑惑啊，莱曼切克小姐

2136
01:41:49,858 --> 01:41:50,838
我这里有封信

2137
01:41:50,918 --> 01:41:53,818
上面说你因为私人原因辞职的

2138
01:41:53,889 --> 01:41:55,189
这封信是谎言

2139
01:41:57,529 --> 01:41:59,119
我是被辞退的，

2140
01:41:59,200 --> 01:42:01,130
这样公司就可以说我不在那里干了

2141
01:42:01,200 --> 01:42:03,070
你可以向法庭指出...

2142
01:42:03,140 --> 01:42:05,000
是谁让你写了这封信的吗？

2143
01:42:13,212 --> 01:42:14,802
杰克.安德霍尔

2144
01:42:16,852 --> 01:42:19,842
他告诉我我必须马上离开，

2145
01:42:19,923 --> 01:42:21,913
而我有两个选择...

2146
01:42:21,993 --> 01:42:24,483
我可以说这是开除，而我就一无所获

2147
01:42:24,563 --> 01:42:26,523
不然我就可以写那封信

2148
01:42:26,594 --> 01:42:28,084
说这是辞职

2149
01:42:28,164 --> 01:42:31,434
而公司会给我1万美元现金作补偿

2150
01:42:31,504 --> 01:42:33,294
让我保持沉默

2151
01:42:33,375 --> 01:42:35,775
而我必须在当时

2152
01:42:35,835 --> 01:42:38,435
当着他的面立刻做选择

2153
01:42:42,946 --> 01:42:44,036
请继续

2154
01:42:45,447 --> 01:42:47,507
我拿了钱...

2155
01:42:48,887 --> 01:42:51,447
然后签了那封信

2156
01:42:51,527 --> 01:42:53,517
保证不和其他人讨论我经手的任何赔付申请

2157
01:42:53,598 --> 01:42:55,588
包括布莱克家的赔付申请

2158
01:42:55,658 --> 01:42:57,688
尤其不能和别人讲布莱克的申请

2159
01:42:57,768 --> 01:43:01,398
那你知道那个申请本应赔付了？

2160
01:43:03,039 --> 01:43:04,529
每个人都知道...

2161
01:43:05,940 --> 01:43:08,840
但是公司在赌

2162
01:43:08,910 --> 01:43:10,400
赌什么？

2163
01:43:11,810 --> 01:43:15,040
赌那个保险人不去咨询律师

2164
01:43:23,592 --> 01:43:25,652
就是说，在那段时间内...

2165
01:43:27,933 --> 01:43:31,303
你是高级赔付审核员？

2166
01:43:31,373 --> 01:43:32,563
是的

2167
01:43:32,634 --> 01:43:34,974
那么在那段期间，你有没有

2168
01:43:35,044 --> 01:43:36,164
得到任何指示

2169
01:43:36,244 --> 01:43:38,714
告诉你该怎么处理这些申请？

2170
01:43:38,774 --> 01:43:41,074
一年之内拒绝一切赔付

2171
01:43:41,145 --> 01:43:43,245
把省下来的钱加起来

2172
01:43:43,315 --> 01:43:45,345
再扣掉庭外和解需要花掉的部分

2173
01:43:45,415 --> 01:43:47,415
还有一大笔剩下来

2174
01:43:47,486 --> 01:43:48,686
我可以靠近书记员吗，法官大人？

2175
01:43:48,756 --> 01:43:49,726
要求批准

2176
01:43:49,796 --> 01:43:50,686
谢谢

2177
01:43:50,756 --> 01:43:52,346
请给我第6号证物

2178
01:43:53,627 --> 01:43:54,997
这是被标记为

2179
01:43:55,067 --> 01:43:57,227
第6号证物的文件

2180
01:43:58,267 --> 01:43:59,897
你认得它吗？

2181
01:43:59,968 --> 01:44:02,908
是的，这是GB保险公司的赔付手册

2182
01:44:04,338 --> 01:44:05,398
你可以帮我个忙

2183
01:44:05,478 --> 01:44:09,468
翻到第U项吗？

2184
01:44:13,120 --> 01:44:15,920
这里没有第U项

2185
01:44:15,990 --> 01:44:17,620
但是你在做高级赔付审核员的时候

2186
01:44:17,690 --> 01:44:19,550
你记得那里有第U项？

2187
01:44:20,661 --> 01:44:22,461
是的，我记得

2188
01:44:22,531 --> 01:44:25,361
第U项是高级审核员手册中的...

2189
01:44:25,431 --> 01:44:27,961
一项行政备忘录

2190
01:44:28,032 --> 01:44:29,592
谢谢你

2191
01:44:29,672 --> 01:44:31,532
法官大人

2192
01:44:31,602 --> 01:44:35,042
这是杰琪.莱曼切克真正的高级赔付手册

2193
01:44:35,113 --> 01:44:36,843
在里面，有一项行政备忘录

2194
01:44:36,913 --> 01:44:38,213
被列为第U项

2195
01:44:38,283 --> 01:44:39,273
我请求上前盘问...

2196
01:44:39,343 --> 01:44:41,183
反对，法官大人，我们可以上前谈谈吗？

2197
01:44:43,254 --> 01:44:46,384
法官大人，我们得到的文件不全

2198
01:44:46,455 --> 01:44:47,685
法官大人，这份文件中的...

2199
01:44:47,755 --> 01:44:49,885
行政备忘录部分是非法窃取的

2200
01:44:49,965 --> 01:44:51,015
不能作为证据

2201
01:44:51,095 --> 01:44:52,025
不能作为证据？

2202
01:44:52,095 --> 01:44:53,285
立奥，你能够证明吗？

2203
01:44:53,366 --> 01:44:55,696
法官大人，我请求你要求我的这位同事

2204
01:44:55,766 --> 01:44:56,756
不要读它

2205
01:44:56,836 --> 01:44:59,066
也不要对任何嫌疑...

2206
01:44:59,136 --> 01:45:02,666
我不理解这为什么不能作为证据

2207
01:45:02,737 --> 01:45:04,367
它不是在正确的时间内获得的

2208
01:45:04,447 --> 01:45:05,567
我们也不知道它到底是怎么获得的

2209
01:45:05,647 --> 01:45:07,477
我昨天晚上才知道这个东西存在的

2210
01:45:07,548 --> 01:45:09,448
假设我现在不接受这项证物

2211
01:45:09,518 --> 01:45:11,248
你还有什么其他问题要问证人的吗？

2212
01:45:12,148 --> 01:45:13,118
没了，法官大人

2213
01:45:13,188 --> 01:45:15,918
你可以交叉询问了，德拉门德先生

2214
01:45:15,989 --> 01:45:17,549
谢谢你

2215
01:45:31,541 --> 01:45:33,801
莱曼切克小姐...

2216
01:45:35,082 --> 01:45:37,712
请问你最近是否

2217
01:45:37,782 --> 01:45:39,712
入院治疗很多疾病？

2218
01:45:39,782 --> 01:45:41,912
我没有入院

2219
01:45:41,983 --> 01:45:45,153
我只是喝了很多酒，精神压抑

2220
01:45:45,223 --> 01:45:48,663
我是自愿去医院检查的

2221
01:45:48,724 --> 01:45:50,564
我的费用本来应该由GB保险公司...

2222
01:45:50,634 --> 01:45:52,824
进行赔付的，我加入了他们的集团保单

2223
01:45:52,894 --> 01:45:55,334
而他们，当然，把我的赔付申请拒绝了

2224
01:45:55,405 --> 01:45:57,265
这大概就是你为什么在这里的原因吧？

2225
01:45:57,335 --> 01:45:59,705
因为你恨GB保险公司？

2226
01:45:59,775 --> 01:46:00,745
莱曼切克小姐？

2227
01:46:00,806 --> 01:46:03,406
我恨GB保险公司

2228
01:46:03,476 --> 01:46:06,246
我也恨大多数在那工作的蛀虫

2229
01:46:06,316 --> 01:46:08,546
那当你和拉佛金先生上床的时候

2230
01:46:08,617 --> 01:46:10,177
你是不是觉得他也是蛀虫之一？

2231
01:46:11,217 --> 01:46:13,317
反对

2232
01:46:13,387 --> 01:46:14,587
法官大人

2233
01:46:14,658 --> 01:46:16,918
德拉门德先生可能觉得谈这个很有趣

2234
01:46:16,998 --> 01:46:19,188
但这个与本案无关

2235
01:46:19,258 --> 01:46:21,128
哦，我可不觉得有趣

2236
01:46:21,199 --> 01:46:22,599
反对无效

2237
01:46:22,669 --> 01:46:23,899
让我们看看这能告诉我们什么

2238
01:46:23,969 --> 01:46:28,459
你是否承认和拉佛金先生有特殊关系？

2239
01:46:34,551 --> 01:46:37,171
莱曼切克小姐？

2240
01:46:42,122 --> 01:46:46,722
在我为GB保险公司工作的时候

2241
01:46:46,792 --> 01:46:49,262
如果我和一些高层人员上床

2242
01:46:49,333 --> 01:46:53,663
我就可以得到提升，工资也能上涨

2243
01:46:53,733 --> 01:46:56,633
而我如果不愿意，立刻就会被降职

2244
01:46:56,704 --> 01:46:58,194
莱曼切克小姐

2245
01:46:58,274 --> 01:47:00,104
作为GB保险公司的雇员

2246
01:47:00,174 --> 01:47:02,004
你保证过不把

2247
01:47:02,075 --> 01:47:04,415
这些保密的私人申请的内容透露给任何人

2248
01:47:04,485 --> 01:47:05,915
- 是不是？
- 是的。

2249
01:47:05,985 --> 01:47:07,345
事实上，你还作证说

2250
01:47:07,415 --> 01:47:10,185
你因此索要

2251
01:47:10,256 --> 01:47:13,246
1万美元，是不是？

2252
01:47:13,326 --> 01:47:15,226
我没有索要过

2253
01:47:15,297 --> 01:47:16,817
但是你接受了，是不是？

2254
01:47:16,897 --> 01:47:18,827
你把钱装进自己的口袋

2255
01:47:18,897 --> 01:47:20,127
而你的心中

2256
01:47:20,197 --> 01:47:22,297
却丝毫没有要保守这些秘密的意思

2257
01:47:22,368 --> 01:47:24,928
事实上，你对GB保险公司

2258
01:47:25,008 --> 01:47:27,438
和拉佛金先生切齿痛恨，是不是？

2259
01:47:27,508 --> 01:47:29,438
你不懂吗，他们在剥削我

2260
01:47:29,509 --> 01:47:31,709
只因为我破产，又是单身

2261
01:47:31,779 --> 01:47:33,179
还有两个小孩

2262
01:47:33,249 --> 01:47:34,939
而你告诉他你打算去见他的妻子

2263
01:47:35,020 --> 01:47:35,950
把这事向报纸捅出去

2264
01:47:36,020 --> 01:47:37,080
而那10000美元

2265
01:47:37,150 --> 01:47:39,310
只是勒索，是不是？

2266
01:47:39,390 --> 01:47:42,620
只是你向你痛恨的公司进行的一种勒索!

2267
01:47:42,691 --> 01:47:43,821
是不是？

2268
01:47:43,891 --> 01:47:45,121
不，不是

2269
01:47:45,191 --> 01:47:47,351
事实上，你今天的作证完全是一场谎言

2270
01:47:47,431 --> 01:47:49,021
你偷窃了公司的文件

2271
01:47:49,102 --> 01:47:52,362
还有机密材料，藉以进行讹诈

2272
01:47:52,432 --> 01:47:56,842
恶魔也不过如此吧，莱曼切克小姐？

2273
01:47:56,913 --> 01:48:00,033
我蔑视你!

2274
01:48:00,113 --> 01:48:03,083
我法官大人，我要求将所有这些...

2275
01:48:03,144 --> 01:48:06,314
莱曼切克小姐提供给控方的文件...

2276
01:48:06,384 --> 01:48:08,254
都作为窃取的文件

2277
01:48:08,325 --> 01:48:10,085
禁止列为证物

2278
01:48:18,366 --> 01:48:22,736
根据法庭所看到的事实

2279
01:48:22,807 --> 01:48:24,497
这些文件不予接受

2280
01:48:24,577 --> 01:48:25,627
喔

2281
01:48:29,508 --> 01:48:31,038
我没有问题了

2282
01:48:31,118 --> 01:48:32,708
谢谢你，德拉门德先生

2283
01:48:32,778 --> 01:48:34,448
贝勒先生

2284
01:48:36,859 --> 01:48:38,949
莱曼切克小姐，你可以下来了

2285
01:48:46,530 --> 01:48:48,230
对不起

2286
01:49:03,153 --> 01:49:04,413
你好

2287
01:49:04,483 --> 01:49:06,183
这是迭克.席夫列特

2288
01:49:06,253 --> 01:49:08,243
我要找大个儿雷诺，你能帮我接一下吗？

2289
01:49:08,323 --> 01:49:10,483
大个儿雷诺？稍等

2290
01:49:10,564 --> 01:49:11,654
好

2291
01:49:16,635 --> 01:49:17,695
你好

2292
01:49:17,765 --> 01:49:19,235
嘿，老板，我是迭克

2293
01:49:19,305 --> 01:49:20,765
哦，嗨，迭克，你怎么样啊？

2294
01:49:20,835 --> 01:49:21,925
很好，你呢？

2295
01:49:22,005 --> 01:49:23,065
哦，我也很好

2296
01:49:23,136 --> 01:49:24,866
你还在这里吗？

2297
01:49:24,946 --> 01:49:26,406
这个嘛，我反正在什么地方

2298
01:49:26,476 --> 01:49:28,566
啊，对

2299
01:49:28,646 --> 01:49:32,806
听着，我有个案子，证据是窃取的

2300
01:49:32,887 --> 01:49:36,017
窃取的证据...

2301
01:49:36,088 --> 01:49:37,818
让我想想看，呃...

2302
01:49:37,888 --> 01:49:40,328
你查一下德.索托的案子

2303
01:49:40,398 --> 01:49:41,358
德.索托？

2304
01:49:41,428 --> 01:49:44,868
卡尔曼.德.索托。记得他吗？

2305
01:49:44,929 --> 01:49:45,829
红宝石俱乐部

2306
01:49:45,899 --> 01:49:47,269
是的

2307
01:49:47,339 --> 01:49:49,269
我在哪能找到这案子？

2308
01:49:49,340 --> 01:49:52,770
呃，大概在92年那卷，

2309
01:49:52,840 --> 01:49:55,240
我想想看，好像是650...

2310
01:49:55,310 --> 01:49:56,680
西南第二上诉庭

2311
01:49:56,751 --> 01:49:58,141
老板，你真救了我的命了

2312
01:49:58,211 --> 01:49:59,151
92年？

2313
01:49:59,211 --> 01:50:00,841
是的，红宝石俱乐部的案子

2314
01:50:00,921 --> 01:50:02,751
92年...

2315
01:50:02,822 --> 01:50:03,842
卡尔曼.德.索托

2316
01:50:03,922 --> 01:50:04,982
卡尔曼.德.索托

2317
01:50:05,052 --> 01:50:05,992
红宝石俱乐部

2318
01:50:06,052 --> 01:50:06,952
红宝石俱乐部

2319
01:50:07,022 --> 01:50:08,282
那是上诉法庭的案子

2320
01:50:08,362 --> 01:50:10,262
对，对，我想起来了

2321
01:50:14,663 --> 01:50:15,633
我找到了

2322
01:50:15,703 --> 01:50:17,763
从布鲁泽那里弄来的

2323
01:50:17,834 --> 01:50:19,774
布鲁泽？我还以为你不知道布鲁泽在哪里呢

2324
01:50:19,834 --> 01:50:22,804
我，我是不知道，但是我有个紧急联络号码

2325
01:50:22,874 --> 01:50:25,034
我给他们打了个电话，他们给我接通了

2326
01:50:25,115 --> 01:50:26,305
鲁迪

2327
01:50:26,375 --> 01:50:28,605
没人比布鲁泽知道更多关于窃取证据的事情了

2328
01:50:28,685 --> 01:50:30,615
这算他的看家本领了

2329
01:50:30,686 --> 01:50:31,776
喂？

2330
01:50:31,856 --> 01:50:33,146
鲁迪.贝勒先生？

2331
01:50:33,216 --> 01:50:34,116
是的

2332
01:50:34,186 --> 01:50:35,886
我是谢尔比郡地区检察官

2333
01:50:35,956 --> 01:50:37,386
你今晚能来下我们法庭吗？

2334
01:50:37,457 --> 01:50:39,557
我想和你谈谈凯莉.莱克的案子

2335
01:50:39,627 --> 01:50:40,557
当然

2336
01:50:40,627 --> 01:50:42,067
来6号房间

2337
01:50:43,868 --> 01:50:46,528
我们有办法对付这群混蛋啦

2338
01:50:46,598 --> 01:50:49,568
在卡尔曼.德.索托告红宝石俱乐部的案子里

2339
01:50:49,638 --> 01:50:51,578
有好多的窃取而来的证据

2340
01:50:51,639 --> 01:50:53,629
布鲁泽自己搞的这个案子...

2341
01:50:53,709 --> 01:50:56,179
现在不行，我得走了

2342
01:50:56,249 --> 01:50:57,479
你要走？

2343
01:50:57,550 --> 01:50:59,980
GB保险公司的首席执行官明天就要来了

2344
01:51:00,050 --> 01:51:03,280
我说，我们得给他好看的

2345
01:51:14,972 --> 01:51:16,462
祝你好运，凯莉

2346
01:51:16,532 --> 01:51:18,002
谢谢

2347
01:51:22,543 --> 01:51:24,443
我和地区检察官谈过了

2348
01:51:24,514 --> 01:51:26,414
她不打算起诉

2349
01:51:26,484 --> 01:51:29,574
她认为你没有罪

2350
01:51:29,654 --> 01:51:31,744
这只是自卫，凯莉

2351
01:51:43,296 --> 01:51:45,566
控方打算传召...

2352
01:51:45,637 --> 01:51:48,467
维尔佛莱德.奇利上庭作证

2353
01:51:48,537 --> 01:51:50,167
请举起你的右手

2354
01:51:50,247 --> 01:51:52,147
你愿意庄严宣誓

2355
01:51:52,208 --> 01:51:53,368
你在此所作证供

2356
01:51:53,448 --> 01:51:54,808
全是事实，是全部事实

2357
01:51:54,878 --> 01:51:56,778
并无丝毫虚言吗？

2358
01:51:56,848 --> 01:51:57,778
是的

2359
01:51:57,849 --> 01:51:59,619
你可以作证了

2360
01:52:01,959 --> 01:52:04,249
请你告诉法庭你的名字，以便记录

2361
01:52:04,330 --> 01:52:05,520
维尔佛莱德.奇利

2362
01:52:05,590 --> 01:52:07,460
我可以上前盘问证人吗？

2363
01:52:07,530 --> 01:52:08,960
要求批准

2364
01:52:09,030 --> 01:52:12,360
奇利先生，呃，在这个GB保险公司的册子里面，

2365
01:52:12,431 --> 01:52:14,231
这个是不是你？这是你的名字吗？

2366
01:52:14,301 --> 01:52:15,231
是的

2367
01:52:15,301 --> 01:52:17,601
那这些缩写指什么？

2368
01:52:17,672 --> 01:52:19,072
你是说CEO？

2369
01:52:19,142 --> 01:52:22,082
是的，它们指什么，CEO，指什么？

2370
01:52:22,142 --> 01:52:23,582
首席执行官

2371
01:52:23,652 --> 01:52:25,882
首席执行官。谢谢你，谢谢你

2372
01:52:25,953 --> 01:52:28,713
就是说，你是...你是头儿了

2373
01:52:28,783 --> 01:52:31,583
你是最重要的人，你说了算

2374
01:52:31,654 --> 01:52:33,314
钱都由你控制了

2375
01:52:33,394 --> 01:52:34,794
是的，你可以这么说

2376
01:52:34,864 --> 01:52:36,054
是吗，好，谢谢你

2377
01:52:36,124 --> 01:52:38,964
呃，法官大人，我希望把询问奇利先生

2378
01:52:39,035 --> 01:52:43,835
的权利交给我的合伙人鲁迪.贝勒

2379
01:52:43,905 --> 01:52:45,495
你疯了吗？

2380
01:52:45,576 --> 01:52:46,666
你还没有执照呢

2381
01:52:46,736 --> 01:52:48,676
我有的选吗，你在干嘛？

2382
01:52:48,746 --> 01:52:50,146
你来晚了

2383
01:52:50,216 --> 01:52:51,506
早上好，法官大人

2384
01:52:51,577 --> 01:52:53,547
对不起，我来晚了

2385
01:52:53,617 --> 01:52:55,917
请求上前盘问证人，法官大人？

2386
01:52:55,987 --> 01:52:57,717
要求批准

2387
01:52:57,788 --> 01:52:59,418
这是杰琪.莱曼切克...

2388
01:52:59,488 --> 01:53:01,418
给我的赔付手册

2389
01:53:01,488 --> 01:53:02,418
反对，法官大人

2390
01:53:02,488 --> 01:53:03,978
窃取的文件不能接受

2391
01:53:04,058 --> 01:53:04,988
你已经对此做过判决了

2392
01:53:05,059 --> 01:53:06,219
反对有效

2393
01:53:06,299 --> 01:53:07,629
我们可以上前谈谈吗？

2394
01:53:12,000 --> 01:53:14,060
法官大人，我认为这个事情已经解决了？

2395
01:53:14,140 --> 01:53:16,070
法官大人，我今天早上

2396
01:53:16,140 --> 01:53:17,660
找到了一个使用

2397
01:53:17,740 --> 01:53:19,070
这样证据的案例

2398
01:53:19,141 --> 01:53:20,511
是什么？

2399
01:53:20,581 --> 01:53:23,071
如果你看看这个判决

2400
01:53:23,151 --> 01:53:25,711
是关于卡尔曼.德.索托和红宝石俱乐部的案子的

2401
01:53:25,782 --> 01:53:28,182
我给您和德拉门德先生都复印了一份

2402
01:53:28,252 --> 01:53:31,222
585号案，西南部第二上诉法庭，第431页

2403
01:53:31,292 --> 01:53:34,052
由布鲁泽，哦不，由小莱曼.斯通事务所辩护

2404
01:53:34,133 --> 01:53:35,353
这里面清楚的表明

2405
01:53:35,433 --> 01:53:37,453
窃取而得的文件，事实上，是可以接受的

2406
01:53:37,533 --> 01:53:40,223
只要该文件并非由律师窃取的即可

2407
01:53:40,304 --> 01:53:42,964
这样啊，嗯，根据这一注释

2408
01:53:43,034 --> 01:53:45,634
你的反对无效

2409
01:53:45,705 --> 01:53:47,005
对不起，立奥

2410
01:53:47,075 --> 01:53:50,045
好吧，你可以这么做，法官大人。

2411
01:53:50,115 --> 01:53:51,845
但请记录下来我对此强烈反对

2412
01:53:51,916 --> 01:53:53,046
你的反对已经被记录了

2413
01:53:53,116 --> 01:53:55,056
我可以上前盘问了吗？

2414
01:53:55,116 --> 01:53:56,386
可以

2415
01:53:58,286 --> 01:54:00,776
呃... 法官大人...

2416
01:54:00,857 --> 01:54:03,327
法官大人，我很抱歉

2417
01:54:03,397 --> 01:54:05,697
对不起，耽误你的时间了，奇利先生

2418
01:54:05,768 --> 01:54:07,628
我们得谈谈这个

2419
01:54:07,698 --> 01:54:09,728
GB保险公司的手册

2420
01:54:09,798 --> 01:54:13,968
这是一份完整的手册吗，先生？

2421
01:54:15,839 --> 01:54:16,969
是的

2422
01:54:17,049 --> 01:54:19,139
其中有第U项吗？

2423
01:54:23,580 --> 01:54:27,210
是的

2424
01:54:27,291 --> 01:54:29,311
那我们谈谈这个神秘的第U项好不好

2425
01:54:29,391 --> 01:54:31,361
你干嘛不向陪审团解释一下呢？

2426
01:54:31,431 --> 01:54:33,221
我们来看一看

2427
01:54:34,402 --> 01:54:36,302
请读出第三段

2428
01:54:40,443 --> 01:54:43,733
“赔付审核员应该在接到任何申请三日内

2429
01:54:43,813 --> 01:54:46,373
毫无例外的，

2430
01:54:46,444 --> 01:54:47,704
拒绝该赔付。”

2431
01:54:47,784 --> 01:54:50,644
你怎么解释这第U项？

2432
01:54:50,714 --> 01:54:53,814
这个嘛，我们有时会收到一些

2433
01:54:53,885 --> 01:54:57,375
无意义和虚假的赔付要求

2434
01:54:57,455 --> 01:54:59,855
因此我们用这一方法

2435
01:54:59,926 --> 01:55:02,196
使我们可以更关注那些

2436
01:55:02,266 --> 01:55:04,456
真实和确实需要的赔付

2437
01:55:05,937 --> 01:55:08,597
奇利先生，你真的认为本法庭...

2438
01:55:08,667 --> 01:55:10,497
会相信你这个解释吗？

2439
01:55:10,577 --> 01:55:12,167
这一段不过是段

2440
01:55:12,237 --> 01:55:15,177
内部处理指示而已

2441
01:55:15,248 --> 01:55:16,708
内部处理指示？

2442
01:55:16,778 --> 01:55:19,008
不，奇利先生，事实并非如此

2443
01:55:19,078 --> 01:55:20,848
第U项所作的远不只如此而已

2444
01:55:20,919 --> 01:55:22,679
孩子，我可不这么认为

2445
01:55:22,749 --> 01:55:25,019
奇利先生，这里面是不是很明确的指出了

2446
01:55:25,089 --> 01:55:27,559
申请应该怎么被传递，转送，再传递

2447
01:55:27,630 --> 01:55:29,620
用尽各种方法，就是拒绝赔付？

2448
01:55:29,690 --> 01:55:30,750
我没看到

2449
01:55:30,830 --> 01:55:32,190
法官大人，请求上前盘问证人？

2450
01:55:32,260 --> 01:55:33,560
要求批准

2451
01:55:37,801 --> 01:55:40,071
那么，奇利先生

2452
01:55:40,142 --> 01:55:44,102
在1995年，GB保险公司有多少份...

2453
01:55:44,182 --> 01:55:47,912
保单生效？

2454
01:55:51,023 --> 01:55:53,283
我不知道

2455
01:55:53,354 --> 01:55:54,844
我们来看看怎么样

2456
01:55:56,754 --> 01:56:00,884
你觉得9万8千份这个数字是否正确，多还是少了？

2457
01:56:02,365 --> 01:56:03,885
差不多吧

2458
01:56:03,965 --> 01:56:05,765
可能就是这么多，是的

2459
01:56:05,835 --> 01:56:06,995
谢谢

2460
01:56:07,066 --> 01:56:10,666
好，在那么这些保单中，有多少人提出了赔付申请？

2461
01:56:12,506 --> 01:56:15,276
这个，这我不知道

2462
01:56:15,347 --> 01:56:18,677
那你觉得1万1千4这个数字对不对

2463
01:56:18,747 --> 01:56:20,377
多还是少了？

2464
01:56:20,458 --> 01:56:22,248
可能对吧，

2465
01:56:22,318 --> 01:56:26,118
但是我得查证一下才行

2466
01:56:26,188 --> 01:56:27,718
哦，我懂了，就是说我要的信息

2467
01:56:27,799 --> 01:56:29,159
都在那本书里了？

2468
01:56:29,229 --> 01:56:30,129
是的

2469
01:56:30,199 --> 01:56:31,529
那你可不可以告诉陪审团

2470
01:56:31,599 --> 01:56:33,969
在这1万1千多份赔付申请中

2471
01:56:34,040 --> 01:56:35,470
有多少被拒绝了？

2472
01:56:38,010 --> 01:56:40,940
我可能找不到，这需要更多时间

2473
01:56:41,011 --> 01:56:43,071
奇利先生，你已经找了两个月了

2474
01:56:43,141 --> 01:56:44,981
现在回答这个问题！

2475
01:56:49,052 --> 01:56:51,352
这个，我...

2476
01:56:52,352 --> 01:56:54,122
我觉得...

2477
01:56:56,323 --> 01:56:57,593
呃...

2478
01:56:59,333 --> 01:57:01,963
9141

2479
01:57:03,464 --> 01:57:06,234
11462份申请

2480
01:57:08,405 --> 01:57:11,395
9141份被据

2481
01:57:14,146 --> 01:57:16,676
法官大人...

2482
01:57:16,746 --> 01:57:18,686
我还有一份文件

2483
01:57:18,756 --> 01:57:20,016
这是GB保险公司医学委员会

2484
01:57:20,086 --> 01:57:23,076
提供的一份报告

2485
01:57:24,797 --> 01:57:27,257
依照先例，请求将此文件...

2486
01:57:27,328 --> 01:57:29,198
呈送给奇利先生

2487
01:57:29,268 --> 01:57:31,128
反对，同样理由，法官大人

2488
01:57:32,638 --> 01:57:34,658
反对无效，不过记录下这一反对

2489
01:57:34,739 --> 01:57:36,929
谢谢你

2490
01:57:43,050 --> 01:57:44,780
呃，奇利先生，这是GB保险公司

2491
01:57:44,850 --> 01:57:47,940
自己的医学委员会写的报告

2492
01:57:48,021 --> 01:57:49,421
而你是这个委员会的主席

2493
01:57:50,661 --> 01:57:53,561
请你从第18行读起好吗？

2494
01:57:58,202 --> 01:58:00,222
“由于骨髓移植...

2495
01:58:00,302 --> 01:58:02,822
已经成为标准治疗方式...

2496
01:58:02,903 --> 01:58:06,133
GB公司投资骨髓临床治疗的

2497
01:58:06,203 --> 01:58:08,673
决定将是有利可图的”

2498
01:58:10,444 --> 01:58:11,704
请求上前盘问证人，法官大人？

2499
01:58:11,784 --> 01:58:13,974
要求批准

2500
01:58:14,054 --> 01:58:16,414
谢谢

2501
01:58:16,485 --> 01:58:18,685
你看到了

2502
01:58:18,755 --> 01:58:20,155
再读大声一点

2503
01:58:20,225 --> 01:58:22,555
反对，法官大人，无意义重复

2504
01:58:22,626 --> 01:58:24,686
反对无效

2505
01:58:24,766 --> 01:58:27,756
我希望陪审团成员都能听到这些

2506
01:58:29,337 --> 01:58:31,237
“由于骨髓移植...

2507
01:58:31,307 --> 01:58:33,767
已经成为标准治疗方式...

2508
01:58:33,837 --> 01:58:36,467
GB公司投资骨髓临床治疗的

2509
01:58:36,538 --> 01:58:39,478
决定将是有利可图的”

2510
01:58:39,548 --> 01:58:41,208
有利可图的

2511
01:58:41,279 --> 01:58:44,719
这就是GB保险公司的方针，是吧？

2512
01:58:45,949 --> 01:58:48,079
我问完证人了

2513
01:58:48,160 --> 01:58:50,020
德拉门德先生

2514
01:58:50,090 --> 01:58:53,060
没有问题。但是我们坚持我们的反对意见

2515
01:58:53,130 --> 01:58:55,250
你可以下去了，奇利先生

2516
01:58:55,331 --> 01:58:57,561
谢谢

2517
01:59:15,924 --> 01:59:16,884
鲁迪

2518
01:59:16,954 --> 01:59:17,854
嗯？

2519
01:59:17,924 --> 01:59:20,184
快点，起床了，快点，到时间了

2520
01:59:20,254 --> 01:59:22,984
快点，快点

2521
01:59:23,065 --> 01:59:28,155
在一个案件中赔付1000万美元意味着什么呢？

2522
01:59:30,266 --> 01:59:33,676
所有的保险金都将用于赔付

2523
01:59:33,736 --> 01:59:35,636
而这个将为....

2524
01:59:35,707 --> 01:59:39,337
政府控制医疗保险铺平道路

2525
01:59:40,717 --> 01:59:44,307
陪审团成员们，你们的职责重大

2526
01:59:44,388 --> 01:59:46,908
我希望你们明智

2527
01:59:46,988 --> 01:59:48,858
谨慎

2528
01:59:48,929 --> 01:59:50,989
公正的作出判决

2529
01:59:53,559 --> 01:59:55,829
公正的作出判决

2530
01:59:57,500 --> 01:59:59,970
谢谢你，法官大人。

2531
02:00:01,310 --> 02:00:04,140
贝勒先生，你可以作结案陈词了

2532
02:00:04,211 --> 02:00:06,611
谢谢你，法官大人

2533
02:00:06,681 --> 02:00:09,171
陪审团的女士们先生们，

2534
02:00:12,082 --> 02:00:14,712
每当我想起丹尼.雷.布莱克...

2535
02:00:16,753 --> 02:00:20,193
呼出最后一口气，慢慢死去的时候...

2536
02:00:21,833 --> 02:00:24,953
我都极度讨厌自己...

2537
02:00:25,034 --> 02:00:27,004
我学了这么多法律

2538
02:00:27,064 --> 02:00:28,664
但是却无法挽救他

2539
02:00:28,734 --> 02:00:32,934
我觉得我是个不合格的律师

2540
02:00:33,005 --> 02:00:34,945
我没有资格作这个结案陈词

2541
02:00:38,946 --> 02:00:42,506
所以，我希望...由丹尼雷代替我做这个陈词

2542
02:00:42,587 --> 02:00:44,577
请展示16号证物

2543
02:00:49,928 --> 02:00:52,898
我现在重110磅（约100斤）

2544
02:00:52,968 --> 02:00:57,398
而11个月前，我还重160磅（约145斤）

2545
02:00:59,669 --> 02:01:02,839
我的白血病很长时间前就确诊了

2546
02:01:04,950 --> 02:01:07,210
我在医院治疗的时候

2547
02:01:07,280 --> 02:01:08,540
医生发现

2548
02:01:08,620 --> 02:01:10,240
唯一可以挽救我的生命的

2549
02:01:10,321 --> 02:01:12,951
就是骨髓移植

2550
02:01:14,521 --> 02:01:19,521
但是我从医院出院了

2551
02:01:19,592 --> 02:01:22,492
以为我们家支付

2552
02:01:22,562 --> 02:01:24,552
不起这笔费用

2553
02:01:24,633 --> 02:01:25,693
你们为什么要这么做？

2554
02:01:28,903 --> 02:01:32,503
GB保险公司拒绝了我们的赔付要求

2555
02:01:32,574 --> 02:01:33,474
过来吧

2556
02:01:33,544 --> 02:01:34,634
没事

2557
02:01:37,755 --> 02:01:40,345
如果我能接受骨髓移植

2558
02:01:40,415 --> 02:01:43,515
我有90％的机会存活下来

2559
02:01:47,026 --> 02:01:49,556
我希望你们和我一样...

2560
02:01:49,627 --> 02:01:51,527
感到震惊

2561
02:01:51,597 --> 02:01:54,497
在如被告这样一家及其富有的大公司

2562
02:01:54,567 --> 02:01:55,627
的长期运作中

2563
02:01:55,707 --> 02:01:58,167
他们一直从低收入家庭中揽钱

2564
02:01:58,238 --> 02:01:59,468
然后拒绝赔付合理的赔付要求

2565
02:01:59,538 --> 02:02:02,978
把那些钱据为己有

2566
02:02:03,049 --> 02:02:06,349
怪不得他们付了那么多钱给他们的律师

2567
02:02:06,419 --> 02:02:09,149
议院游说团和公共关系网...

2568
02:02:09,219 --> 02:02:12,159
要让我们相信我们需要一个改革...

2569
02:02:12,220 --> 02:02:15,620
停止惩罚性的赔偿

2570
02:02:23,402 --> 02:02:26,232
在此，我恳求陪审团...

2571
02:02:27,612 --> 02:02:28,972
作出...

2572
02:02:31,383 --> 02:02:33,783
你们内心深处...

2573
02:02:33,843 --> 02:02:35,543
认为正确的决定

2574
02:02:35,613 --> 02:02:38,783
如果你们不对GB保险公司作出惩罚，

2575
02:02:40,154 --> 02:02:42,784
你们可能就是下一个受害者

2576
02:02:46,895 --> 02:02:48,365
我说完了

2577
02:03:08,488 --> 02:03:10,048
别紧张

2578
02:03:30,181 --> 02:03:32,411
陪审团有判决了吗？

2579
02:03:32,482 --> 02:03:34,142
是的，我们已经作出判决了，法官大人

2580
02:03:34,212 --> 02:03:35,912
你们已经按照我的指示

2581
02:03:35,982 --> 02:03:37,382
将判决写在纸上了吗？

2582
02:03:37,453 --> 02:03:39,013
是的

2583
02:03:39,093 --> 02:03:40,683
请读出你们的判决

2584
02:03:42,323 --> 02:03:46,593
“我们，陪审团，认为原告方

2585
02:03:46,664 --> 02:03:49,154
应被赔偿其实际损失

2586
02:03:49,234 --> 02:03:53,334
15万美元”

2587
02:03:56,675 --> 02:04:02,015
“并且，我们，陪审团认为被告方

2588
02:04:02,086 --> 02:04:04,546
应付与原告方惩罚性赔偿

2589
02:04:04,617 --> 02:04:08,817
5000万美元”

2590
02:04:24,339 --> 02:04:26,139
全体起立

2591
02:04:27,440 --> 02:04:28,930
这些人试图隐瞒事实

2592
02:04:29,010 --> 02:04:31,680
逃避责任...

2593
02:04:31,751 --> 02:04:32,941
最后他们都被抓到了

2594
02:04:33,021 --> 02:04:34,641
关于这个故事的另一个有趣的事实是...

2595
02:04:34,721 --> 02:04:36,911
这是这个案子的主要律师鲁迪.贝勒

2596
02:04:36,991 --> 02:04:38,611
打的第一场案子

2597
02:04:38,692 --> 02:04:39,992
不管从什么人的观点来看

2598
02:04:40,062 --> 02:04:41,992
这都是个惊人的判决

2599
02:04:42,062 --> 02:04:43,462
这肯定是最高的...

2600
02:04:43,532 --> 02:04:45,892
呵呵，我的小帮工干的不错嘛

2601
02:04:45,963 --> 02:04:47,763
其实，这没什么难的

2602
02:04:47,833 --> 02:04:49,923
我们只是有个很好的陪审团

2603
02:04:50,003 --> 02:04:52,063
然后事实纷纷出现

2604
02:04:52,144 --> 02:04:53,664
5000万美元的惩罚性赔偿

2605
02:04:53,744 --> 02:04:55,004
5000万美元？

2606
02:04:55,074 --> 02:04:56,044
是的

2607
02:04:56,114 --> 02:04:57,664
你能挣多少？

2608
02:04:57,744 --> 02:04:59,774
你觉得现在有钱了，是吧？

2609
02:04:59,845 --> 02:05:01,335
对不起，我不是这个意思

2610
02:05:01,415 --> 02:05:02,645
我知道你不是这个意思的

2611
02:05:02,715 --> 02:05:06,485
我们能拿到其中1/3，不过钱还没有到帐

2612
02:05:07,456 --> 02:05:08,886
待会见

2613
02:05:08,956 --> 02:05:10,896
你就把这女孩这么一个人扔在这？

2614
02:05:10,956 --> 02:05:12,446
不会很长的

2615
02:05:13,997 --> 02:05:16,727
鲁迪。鲁迪

2616
02:05:16,797 --> 02:05:20,097
我打算把你列入我的遗产名单

2617
02:05:20,168 --> 02:05:21,398
真的？你真是...

2618
02:05:21,478 --> 02:05:24,598
你真是对我太好了，伯迪夫人

2619
02:05:28,579 --> 02:05:30,349
热点新闻

2620
02:05:30,419 --> 02:05:34,979
维尔佛莱德.奇利先生，GB保险公司首席执行官

2621
02:05:35,060 --> 02:05:37,520
昨天下午在（纽约）约翰肯尼迪国际机场被拘留

2622
02:05:37,590 --> 02:05:40,580
他当时和她夫人一起

2623
02:05:40,661 --> 02:05:42,061
试图登时一架去（伦敦）希思罗机场的飞机

2624
02:05:42,131 --> 02:05:44,961
他们说自己只是去度一个短假

2625
02:05:45,031 --> 02:05:47,161
但是他们却无法

2626
02:05:47,232 --> 02:05:48,332
举出任何一家..."

2627
02:05:48,402 --> 02:05:49,662
他们打算去的欧洲旅店的名字

2628
02:05:49,742 --> 02:05:52,332
今天下午5时，GB保险公司

2629
02:05:52,413 --> 02:05:54,343
依据破产法，

2630
02:05:54,413 --> 02:05:55,903
在克立夫兰联邦法庭提出破产保护

2631
02:05:55,983 --> 02:05:58,543
多个州政府正在调查GB保险公司

2632
02:05:58,613 --> 02:06:01,273
还有很多连锁赔偿案已经被提起诉讼

2633
02:06:03,854 --> 02:06:05,284
喂？

2634
02:06:05,354 --> 02:06:08,324
鲁迪，我是立奥.德拉门德

2635
02:06:08,395 --> 02:06:12,925
看来整个公司已经被洗劫一空了

2636
02:06:12,996 --> 02:06:14,966
我很抱歉，鲁迪

2637
02:06:15,036 --> 02:06:18,296
我希望你能拿到你的那份

2638
02:06:18,376 --> 02:06:21,276
我想让你知道

2639
02:06:21,347 --> 02:06:23,397
这场官司每个人都没赚到便宜

2640
02:06:24,647 --> 02:06:26,407
谢谢你，德拉门德先生

2641
02:06:27,718 --> 02:06:29,948
GB就像个坏了的赌博机...

2642
02:06:30,018 --> 02:06:32,548
从不付钱

2643
02:06:32,619 --> 02:06:34,589
我们要是拿那17万5就好啦

2644
02:06:34,659 --> 02:06:37,679
我们到底在想什么？

2645
02:06:41,800 --> 02:06:43,790
这事简直乱套了，是不是？

2646
02:06:43,870 --> 02:06:46,530
这种法律...

2647
02:06:46,601 --> 02:06:47,831
美国的每个律师

2648
02:06:47,911 --> 02:06:50,241
现在都在谈论我

2649
02:06:50,311 --> 02:06:51,571
但这些却没有让我觉得

2650
02:06:51,641 --> 02:06:54,771
自己已经是法律事业的一员了

2651
02:06:54,852 --> 02:06:56,752
如果我全身心的投入工作

2652
02:06:56,812 --> 02:06:59,252
我可能能继续从事法律工作

2653
02:06:59,323 --> 02:07:01,253
但是这样我没法照顾凯莉

2654
02:07:01,323 --> 02:07:05,853
而她现在需要非常多的照顾

2655
02:07:07,224 --> 02:07:10,194
所以尽管我还是热爱法律，我会永远热爱它

2656
02:07:10,264 --> 02:07:12,564
也许我可以在哪里教法律

2657
02:07:12,634 --> 02:07:15,864
而不是做律师

2658
02:07:15,935 --> 02:07:18,065
我需要时间慢慢搞定这事

2659
02:07:18,145 --> 02:07:20,665
GB保险公司现在一文不名了

2660
02:07:22,246 --> 02:07:24,876
这曾经好像天上掉下的馅饼

2661
02:07:24,946 --> 02:07:26,536
但是现在我们一分钱也没捞到

2662
02:07:26,617 --> 02:07:28,707
天啊

2663
02:07:28,787 --> 02:07:31,087
你把他们弄出局了

2664
02:07:31,157 --> 02:07:34,127
田纳西孟菲斯的一个小女人

2665
02:07:34,188 --> 02:07:37,358
就把这群混蛋弄破产了？

2666
02:07:37,428 --> 02:07:38,418
哦

2667
02:07:40,129 --> 02:07:42,759
我明天会去丹尼雷的墓前，

2668
02:07:42,839 --> 02:07:44,929
告诉他这个事情

2669
02:07:50,910 --> 02:07:53,570
哦，凯莉和我大概

2670
02:07:53,651 --> 02:07:57,741
明天一早就要走了

2671
02:07:57,821 --> 02:07:59,841
我们创造了历史，鲁迪

2672
02:08:01,322 --> 02:08:02,692
你知道吗？

2673
02:08:05,362 --> 02:08:07,332
我们可以吹一辈子了

2674
02:08:09,403 --> 02:08:10,833
再见

2675
02:08:11,973 --> 02:08:13,693
毫无疑问我已经达到了事业的顶峰

2676
02:08:13,774 --> 02:08:15,834
事实上，由于这个案子把我送到顶峰

2677
02:08:15,904 --> 02:08:18,704
这之后我不管怎么走都是下坡路了

2678
02:08:18,774 --> 02:08:19,804
再见

2679
02:08:19,874 --> 02:08:22,214
我以后碰到的每一个客户都会期盼这种...

2680
02:08:22,285 --> 02:08:25,405
这种奇迹，一点都不能少

2681
02:08:25,485 --> 02:08:27,915
而也许我确实可以给他们这种奇迹

2682
02:08:27,986 --> 02:08:30,586
如果我不在意用什么方法的话

2683
02:08:30,656 --> 02:08:32,916
但是那样的话，有个早上当我醒来，

2684
02:08:32,996 --> 02:08:35,796
我会发现自己已经变成了立奥.德拉门德

2685
02:08:35,867 --> 02:08:40,327
嘿，你通过律考以后给我个电话好吗？

2686
02:08:40,398 --> 02:08:42,628
没问题，轻而易举！

2687
02:08:45,008 --> 02:08:47,878
每个律师，起码是这个案中的每个律师

2688
02:08:47,939 --> 02:08:51,139
都发现自己越过了一条自己不愿意越过的线

2689
02:08:51,209 --> 02:08:52,269
事情就是这样

2690
02:08:52,349 --> 02:08:54,339
如果你一次次越过这条线

2691
02:08:54,420 --> 02:08:56,390
有天你会发现这条线已经不存在了

2692
02:08:57,620 --> 02:09:00,320
不过那时你自己也就只是个律师笑话了

2693
02:09:00,391 --> 02:09:03,421
只是污水里的又一条鲨鱼`)
	if len(lang) != 1 || lang[0] != "zh" {
		t.Errorf("Expect zh but %v", lang)
	}
}
