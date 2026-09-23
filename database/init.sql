CREATE TABLE IF NOT EXISTS users (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(80) NOT NULL,
  major VARCHAR(120) NOT NULL,
  credit_score INT NOT NULL DEFAULT 80,
  credit_level VARCHAR(40) NOT NULL
);

CREATE TABLE IF NOT EXISTS skills (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  title VARCHAR(120) NOT NULL,
  category VARCHAR(40) NOT NULL,
  level_score INT NOT NULL,
  campus VARCHAR(40) NOT NULL,
  description TEXT NOT NULL,
  portfolio VARCHAR(160) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS needs (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  title VARCHAR(120) NOT NULL,
  category VARCHAR(40) NOT NULL,
  campus VARCHAR(40) NOT NULL,
  expect_time VARCHAR(80) NOT NULL,
  budget_type VARCHAR(40) NOT NULL,
  description TEXT NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'open' COMMENT 'open=开放中 closed=已关闭'
);

-- 求助候选名单：响应者按提交先后排队，发布人挑一位约谈
CREATE TABLE IF NOT EXISTS need_responses (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  need_id BIGINT NOT NULL,
  responder VARCHAR(80) NOT NULL,
  help_method VARCHAR(200) NOT NULL COMMENT '能帮的方式',
  time_slots VARCHAR(200) NOT NULL COMMENT '可上门时段，逗号分隔',
  status VARCHAR(20) NOT NULL DEFAULT 'waiting' COMMENT 'waiting=排队中 interviewing=约谈中 withdrawn=已撤回 exited=已退出 closed=轮候结束',
  position INT NOT NULL COMMENT '按提交先后的排队位次',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_need_position (need_id, position)
);

CREATE TABLE IF NOT EXISTS appointments (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  pair_name VARCHAR(120) NOT NULL,
  exchange_time VARCHAR(80) NOT NULL,
  place VARCHAR(120) NOT NULL,
  status VARCHAR(40) NOT NULL,
  agenda TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS reviews (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  from_user VARCHAR(80) NOT NULL,
  to_user VARCHAR(80) NOT NULL,
  rating INT NOT NULL,
  content TEXT NOT NULL
);

INSERT INTO users(name, major, credit_score, credit_level) VALUES
('林澈', '新闻传播 2023', 91, '黄金导师'),
('孟野', '音乐表演 2022', 88, '白银协作者'),
('周芮', '统计学 2021', 93, '黄金导师');

INSERT INTO skills(user_id, title, category, level_score, campus, description, portfolio) VALUES
(1, '毕业照人像摄影', '摄影', 92, '东校区', '提供构图、修图和毕业季跟拍，可交换吉他入门课。', '12组校园人像作品'),
(2, '民谣吉他陪练', '乐器', 81, '西校区', '节奏型、弹唱和舞台经验分享，想找人拍宣传照。', '校园音乐节演出视频'),
(3, 'Python 数据分析', '编程', 88, '中心校区', 'pandas、可视化、论文数据清洗辅导。', '3份课程项目证书');

INSERT INTO needs(user_id, title, category, campus, expect_time, budget_type, description, status) VALUES
(2, '找人帮忙拍乐队宣传照', '摄影', '西校区', '本周六上午', '技能交换', '可交换 3 次吉他课，希望会调色和室外构图。', 'open'),
(1, '想学吉他扫弦入门', '乐器', '东校区', '周三晚', '技能交换', '用摄影课交换吉他基础，希望同校区或线上。', 'open');

INSERT INTO need_responses(need_id, responder, help_method, time_slots, status, position) VALUES
(1, '林澈', '室外人像跟拍 + 精修 9 张，可带补光灯', '周六上午,周日下午', 'interviewing', 1),
(1, '周芮', '现场花絮记录，可带反光板协助打光', '周六上午', 'waiting', 2),
(2, '孟野', '民谣扫弦入门 2 次课，想换摄影课', '周三晚,周六上午', 'waiting', 1);
